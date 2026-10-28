// installer: install, import, uninstall and scan apps.

package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/HycJack/macdesk/kit"
)

// maxAppBytes caps one package's installed size.
const maxAppBytes = 2 << 20

// appsRoot is where installed apps live; tests repoint it.
var appsRoot = "apps"

// installApp validates a package directory and installs it.
func (d *desk) installApp(dir string) error {
	mb, err := os.ReadFile(filepath.Join(dir, "manifest.json"))
	if err != nil {
		return fmt.Errorf("apphost: no manifest: %w", err)
	}
	var m appManifest
	if err := json.Unmarshal(mb, &m); err != nil {
		return fmt.Errorf("apphost: bad manifest: %w", err)
	}
	if err := checkManifest(&m); err != nil {
		return err
	}
	if ia := d.installed[m.ID]; ia != nil && ia.manifest.Version == m.Version {
		return fmt.Errorf("apphost: %s %s already installed", m.ID, m.Version)
	}
	var size int64
	_ = filepath.WalkDir(dir, func(_ string, e os.DirEntry, err error) error {
		if err == nil && !e.IsDir() {
			if info, infoErr := e.Info(); infoErr == nil {
				size += info.Size()
			}
		}
		return nil
	})
	if size > maxAppBytes {
		return fmt.Errorf("apphost: package is %d bytes, over 2 MiB", size)
	}
	src, err := os.ReadFile(filepath.Join(dir, m.entry()))
	if err != nil {
		return fmt.Errorf("apphost: no entry: %w", err)
	}
	if _, err := loadAppSource(m.Name, string(src), m.Permissions,
		filepath.Join(appsRoot, m.ID)); err != nil {
		return err
	}
	dst := filepath.Join(appsRoot, m.ID, m.Version)
	_ = os.RemoveAll(dst)
	if err := copyDir(dir, dst); err != nil {
		return err
	}
	if d.installed == nil {
		d.installed = map[string]*installedApp{}
	}
	d.installed[m.ID] = &installedApp{manifest: m, dir: dst}
	pruneVersions(m.ID, m.Version)
	kit.Log("appstore", "install", "id", m.ID)
	return nil
}

// pruneVersions removes version dirs other than the installed one.
func pruneVersions(id, keep string) {
	entries, err := os.ReadDir(filepath.Join(appsRoot, id))
	if err != nil {
		return
	}
	for _, e := range entries {
		if e.IsDir() && e.Name() != keep {
			_ = os.RemoveAll(filepath.Join(appsRoot, id, e.Name()))
		}
	}
}

// copyDir copies a package directory recursively.
func copyDir(from, to string) error {
	return filepath.WalkDir(from, func(p string, e os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(from, p)
		dst := filepath.Join(to, rel)
		if e.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, b, 0o644)
	})
}

// importApp installs from a .goapp zip or a package directory.
func (d *desk) importApp(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("apphost: empty path")
	}
	low := strings.ToLower(path)
	if strings.HasSuffix(low, ".zip") || strings.HasSuffix(low, ".goapp") {
		return d.importZip(path)
	}
	return d.installApp(path)
}

// importZip extracts a .goapp archive and installs it.
func (d *desk) importZip(path string) error {
	r, err := zip.OpenReader(path)
	if err != nil {
		return fmt.Errorf("apphost: bad archive: %w", err)
	}
	defer r.Close()
	var total int64
	for _, f := range r.File {
		total += int64(f.UncompressedSize64)
	}
	if total > maxAppBytes {
		return fmt.Errorf("apphost: package is %d bytes, over 2 MiB", total)
	}
	tmp, err := os.MkdirTemp("", "goapp-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	for _, f := range r.File {
		dest := filepath.Join(tmp, filepath.Clean("/"+f.Name))
		if !strings.HasPrefix(dest, filepath.Clean(tmp)) {
			return errors.New("apphost: unsafe zip path")
		}
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(dest, 0o755)
			continue
		}
		if err := extractFile(f, dest); err != nil {
			return err
		}
	}
	return d.installApp(tmp)
}

// extractFile copies one archive entry to disk.
func extractFile(f *zip.File, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, rc)
	return err
}

// uninstallApp removes an installed app and closes its window.
func (d *desk) uninstallApp(id string) error {
	ia := d.installed[id]
	if ia == nil {
		return fmt.Errorf("apphost: %q is not installed", id)
	}
	wid := "app:" + id
	if ws := d.wins[wid]; ws != nil {
		ws.open = false
	}
	d.stopApp(id)
	if err := os.RemoveAll(filepath.Join(appsRoot, id)); err != nil {
		return err
	}
	delete(d.installed, id)
	p := d.tr()
	d.pushTrash(ia.manifest.Name)
	d.pushNotif(p.notifRemove+" "+ia.manifest.Name, p.trashHint)
	kit.Log("appstore", "uninstall", "id", id)
	return nil
}

// scanInstalled loads the apps directory into state.
func (d *desk) scanInstalled() {
	d.installed = map[string]*installedApp{}
	ids, err := os.ReadDir(appsRoot)
	if err != nil {
		return
	}
	for _, id := range ids {
		if !id.IsDir() {
			continue
		}
		versions, err := os.ReadDir(filepath.Join(appsRoot, id.Name()))
		if err != nil {
			continue
		}
		versions = slices.DeleteFunc(versions, func(e os.DirEntry) bool {
			return !e.IsDir() || !validVersion(e.Name())
		})
		slices.SortFunc(versions, func(a, b os.DirEntry) int {
			if versionLess(a.Name(), b.Name()) {
				return -1
			}
			return 1
		})
		for v := len(versions) - 1; v >= 0; v-- {
			mb, err := os.ReadFile(filepath.Join(
				appsRoot, id.Name(), versions[v].Name(), "manifest.json"))
			if err != nil {
				continue
			}
			var m appManifest
			if json.Unmarshal(mb, &m) != nil || checkManifest(&m) != nil {
				continue
			}
			d.installed[m.ID] = &installedApp{manifest: m,
				dir: filepath.Join(appsRoot, id.Name(), versions[v].Name())}
			break
		}
	}
}

// versionLess compares dotted numeric versions.
func versionLess(a, b string) bool {
	as, bs := strings.Split(a, "."), strings.Split(b, ".")
	if len(as) != 3 || len(bs) != 3 {
		return a < b
	}
	for i := 0; i < 3; i++ {
		x, _ := strconv.Atoi(as[i])
		y, _ := strconv.Atoi(bs[i])
		if x != y {
			return x < y
		}
	}
	return false
}

// installedByName finds an installed app by display name.
func (d *desk) installedByName(name string) *installedApp {
	for _, ia := range d.installed {
		if ia.manifest.Name == name {
			return ia
		}
	}
	return nil
}
