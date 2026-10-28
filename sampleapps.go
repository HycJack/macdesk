package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed sampleapps
var sampleFS embed.FS

// samplePack is one embedded installable package.
type samplePack struct {
	id, name, version string
	perms             []string
	files             map[string]string
}

// samplePacks are the built-in installable apps.
var samplePacks = loadSamples()

// loadSamples reads the embedded sample packages once.
func loadSamples() []samplePack {
	dirs, err := fs.ReadDir(sampleFS, "sampleapps")
	if err != nil {
		return nil
	}
	var out []samplePack
	for _, e := range dirs {
		if !e.IsDir() {
			continue
		}
		sp := samplePack{id: e.Name(), files: map[string]string{}}
		entries, err := fs.ReadDir(sampleFS, "sampleapps/"+e.Name())
		if err != nil {
			continue
		}
		for _, f := range entries {
			b, err := fs.ReadFile(sampleFS, "sampleapps/"+e.Name()+"/"+f.Name())
			if err != nil {
				continue
			}
			sp.files[f.Name()] = string(b)
			if f.Name() == "manifest.json" {
				var m appManifest
				if json.Unmarshal(b, &m) == nil {
					sp.name, sp.version, sp.perms = m.Name, m.Version, m.Permissions
				}
			}
		}
		out = append(out, sp)
	}
	return out
}

// sampleByName finds an embedded sample by display name.
func sampleByName(name string) *samplePack {
	for i := range samplePacks {
		if samplePacks[i].name == name {
			return &samplePacks[i]
		}
	}
	return nil
}

// sampleByID finds an embedded sample by id.
func sampleByID(id string) *samplePack {
	for i := range samplePacks {
		if samplePacks[i].id == id {
			return &samplePacks[i]
		}
	}
	return nil
}

// materialize writes the package to a temp directory.
func (s samplePack) materialize() (string, error) {
	dir, err := os.MkdirTemp("", "goapp-"+s.id+"-")
	if err != nil {
		return "", err
	}
	for name, content := range s.files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			os.RemoveAll(dir)
			return "", err
		}
	}
	return dir, nil
}

// installSample installs one embedded package.
func (d *desk) installSample(s samplePack) error {
	dir, err := s.materialize()
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	return d.installApp(dir)
}
