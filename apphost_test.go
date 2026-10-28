package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const goodApp = `
package main

import (
	"strconv"

	"appkit"
)

func Open() any { return map[string]any{"n": 0} }

func View(s any) *appkit.Node {
	return appkit.Column(4, appkit.Text("n="+strconv.Itoa(anyInt(s)), "body"))
}

func Action(s any, a map[string]any) any {
	m := s.(map[string]any)
	if a["type"] == "add" {
		m["n"] = anyInt(s) + 1
	}
	return m
}

func anyInt(s any) int { return s.(map[string]any)["n"].(int) }

func Close() {}
`

func TestLoadAppSource(t *testing.T) {
	hooks, err := loadAppSource("demo", goodApp, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	res := hooks.open.Call(nil)
	if m, ok := res[0].Interface().(map[string]any); !ok || m["n"] != 0 {
		t.Errorf("Open gave %v", res[0])
	}
}

func TestLoadAppSourceRejects(t *testing.T) {
	cases := []struct{ name, src, want string }{
		{"syntax", "package main\nfunc Open( {", "does not compile"},
		{"no view", "package main\nfunc Open() any { return nil }", "missing"},
		{"stdlib", "package main\nimport \"os\"\nfunc Open() any { return nil }\n" +
			"func View(s any) *Node { return nil }\n" +
			"func Action(s any, a map[string]any) any { return s }", "does not compile"},
		{"bad view sig", "package main\nimport \"appkit\"\n" +
			"func Open() any { return nil }\n" +
			"func View(s string) *appkit.Node { return nil }\n" +
			"func Action(s any, a map[string]any) any { return s }", "is func(string)"},
	}
	for _, tc := range cases {
		if _, err := loadAppSource(tc.name, tc.src, nil, ""); err == nil {
			t.Errorf("%s: accepted", tc.name)
		} else if !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: %v, want %q", tc.name, err, tc.want)
		}
	}
}

func TestAppInstanceRoundTrip(t *testing.T) {
	inst, err := startAppSource("demo", goodApp, nil, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	inst.doAction(map[string]any{"type": "add"})
	inst.doAction(map[string]any{"type": "add"})
	tree := inst.tree.Load()
	if tree == nil || tree.Children[0].Text != "n=2" {
		t.Fatalf("tree %+v", tree)
	}
}

func TestAppPanicRecovers(t *testing.T) {
	bad := "package main\nimport \"appkit\"\n" +
		"func Open() any { return nil }\n" +
		"func View(s any) *appkit.Node { return appkit.Text(s.(map[string]any)[\"x\"].(string), \"body\") }\n" +
		"func Action(s any, a map[string]any) any { return s }\n"
	inst, err := startAppSource("bad", bad, nil, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	inst.doAction(map[string]any{"type": "tap"})
	if inst.lastErr == "" {
		t.Error("panic did not surface as lastErr")
	}
	if inst.broken != "" {
		t.Error("a crash must not break the app")
	}
}

func TestStoragePermission(t *testing.T) {
	dir := t.TempDir()
	src := "package main\nimport \"appkit\"\n" +
		"func Open() any { appkit.StoragePut(\"k\", \"v\"); return appkit.StorageGet(\"k\") }\n" +
		"func View(s any) *appkit.Node { return appkit.Text(s.(string), \"body\") }\n" +
		"func Action(s any, a map[string]any) any { return s }\n"
	inst, err := startAppSource("store", src, []string{"storage"}, dir)
	if err != nil {
		t.Fatal(err)
	}
	root := inst.tree.Load()
	if root == nil || root.Text != "v" {
		t.Error("storage round trip failed")
	}
	if _, err := os.Stat(filepath.Join(dir, "storage.json")); err != nil {
		t.Error("storage.json not written")
	}
	// Without the permission the symbols do not exist.
	if _, err := startAppSource("nostore", src, nil, t.TempDir()); err == nil {
		t.Error("storage reachable without permission")
	}
}

// isolateApps points the apps root at a temp directory.
func isolateApps(t *testing.T) {
	t.Helper()
	old := appsRoot
	appsRoot = t.TempDir()
	t.Cleanup(func() { appsRoot = old })
}

func writePackage(t *testing.T, dir, id string) string {
	t.Helper()
	root := filepath.Join(dir, id)
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "manifest.json"),
		[]byte(`{"id":"`+id+`","name":"演示","version":"1.0.0","api":1}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte(goodApp), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestInstallAndUninstall(t *testing.T) {
	isolateApps(t)
	d := newDesk()
	src := writePackage(t, t.TempDir(), "demo")
	if err := d.installApp(src); err != nil {
		t.Fatal(err)
	}
	if d.installed["demo"] == nil {
		t.Fatal("not installed")
	}
	if _, err := os.Stat(filepath.Join(appsRoot, "demo", "1.0.0", "main.go")); err != nil {
		t.Fatal(err)
	}
	if d.installedByName("演示") == nil {
		t.Error("lookup by name failed")
	}
	if err := d.uninstallApp("demo"); err != nil {
		t.Fatal(err)
	}
	if d.installed["demo"] != nil {
		t.Error("still installed")
	}
	if _, err := os.Stat(filepath.Join(appsRoot, "demo")); !os.IsNotExist(err) {
		t.Error("directory not removed")
	}
}

func TestInstallRejects(t *testing.T) {
	isolateApps(t)
	d := newDesk()
	tmp := t.TempDir()
	src := writePackage(t, tmp, "demo")
	if err := d.installApp(src); err != nil {
		t.Fatal(err)
	}
	if err := d.installApp(src); err == nil {
		t.Error("same version reinstalled")
	}
	big := writePackage(t, tmp, "big")
	if err := os.WriteFile(filepath.Join(big, "blob.bin"), make([]byte, maxAppBytes+1), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := d.installApp(big); err == nil || !strings.Contains(err.Error(), "2 MiB") {
		t.Errorf("oversize accepted: %v", err)
	}
	bad := writePackage(t, tmp, "bad")
	if err := os.WriteFile(filepath.Join(bad, "main.go"), []byte("package main\nfunc X() {}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := d.installApp(bad); err == nil {
		t.Error("hookless package accepted")
	}
	evil := writePackage(t, tmp, "evil")
	if err := os.WriteFile(filepath.Join(evil, "manifest.json"),
		[]byte(`{"id":"evil","name":"x","version":"..","api":1}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := d.installApp(evil); err == nil {
		t.Error("dotted version accepted")
	}
}

func TestImportZip(t *testing.T) {
	isolateApps(t)
	d := newDesk()
	src := writePackage(t, t.TempDir(), "demo")
	zipPath := filepath.Join(t.TempDir(), "demo.goapp")
	if err := zipDir(src, zipPath); err != nil {
		t.Fatal(err)
	}
	if err := d.importApp(zipPath); err != nil {
		t.Fatal(err)
	}
	if d.installed["demo"] == nil {
		t.Error("zip not installed")
	}
}

func TestImportZipSlip(t *testing.T) {
	isolateApps(t)
	d := newDesk()
	zipPath := filepath.Join(t.TempDir(), "evil.goapp")
	if err := slipZip(zipPath); err != nil {
		t.Fatal(err)
	}
	if err := d.importApp(zipPath); err == nil {
		t.Fatal("zip slip accepted")
	}
	if _, err := os.Stat(filepath.Join("..", "escaped.txt")); !os.IsNotExist(err) {
		t.Error("zip slip wrote outside")
	}
}

// slipZip writes an archive whose entry escapes its root.
func slipZip(dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	w := zip.NewWriter(out)
	f, err := w.Create("../../escaped.txt")
	if err != nil {
		return err
	}
	_, err = f.Write([]byte("x"))
	if err != nil {
		return err
	}
	return w.Close()
}

func TestScanToleratesStorageFile(t *testing.T) {
	isolateApps(t)
	d := newDesk()
	tmp := t.TempDir()
	if err := d.installApp(writePackage(t, tmp, "demo")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(appsRoot, "demo", "storage.json"),
		[]byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	d.scanInstalled()
	if d.installed["demo"] == nil {
		t.Error("app lost after scan with storage file")
	}
}

func TestUpgradePrunesOldVersion(t *testing.T) {
	isolateApps(t)
	d := newDesk()
	tmp := t.TempDir()
	if err := d.installApp(writePackage(t, tmp, "demo")); err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(tmp, "demo")
	if err := os.WriteFile(filepath.Join(root, "manifest.json"),
		[]byte(`{"id":"demo","name":"演示","version":"1.1.0","api":1}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := d.installApp(root); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(appsRoot, "demo", "1.0.0")); !os.IsNotExist(err) {
		t.Error("old version dir kept")
	}
	if d.installed["demo"].manifest.Version != "1.1.0" {
		t.Errorf("version %q", d.installed["demo"].manifest.Version)
	}
}

func TestScanInstalled(t *testing.T) {
	isolateApps(t)
	d := newDesk()
	src := writePackage(t, t.TempDir(), "demo")
	if err := d.installApp(src); err != nil {
		t.Fatal(err)
	}
	fresh := newDesk()
	fresh.scanInstalled()
	if fresh.installed["demo"] == nil {
		t.Error("scan missed the installed app")
	}
}

// zipDir packs a package directory into a zip archive.
func zipDir(src, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	w := zip.NewWriter(out)
	defer w.Close()
	return filepath.WalkDir(src, func(p string, e os.DirEntry, err error) error {
		if err != nil || e.IsDir() {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		f, err := w.Create(rel)
		if err != nil {
			return err
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		_, err = f.Write(b)
		return err
	})
}
