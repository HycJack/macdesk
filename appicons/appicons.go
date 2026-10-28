// Package appicons holds the macOS style icon set: coloured app artwork and
// single colour glyphs, drawn from the symbol sheet of StarKnightt/macos-web
// (MIT). Each icon is an SVG file parsed once and cached.
package appicons

import (
	"embed"
	"errors"
	"io/fs"
	"sort"
	"strings"
	"sync"

	"github.com/egoist/mygo/ui"
)

// svg/ is generated from the symbol sheet; custom/ holds hand-drawn icons and
// the ones whose text had to be redrawn as paths, which win on a shared name.
//
//go:embed svg/*.svg custom/*.svg
var files embed.FS

var (
	mu    sync.Mutex
	cache = map[string]*ui.SVG{}
	dirs  = []string{"custom", "svg"}
)

// read returns the bytes of an icon, preferring custom/ over svg/.
func read(name string) ([]byte, error) {
	for _, dir := range dirs {
		data, err := files.ReadFile(dir + "/" + name + ".svg")
		if err == nil {
			return data, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}
	return nil, fs.ErrNotExist
}

// Get returns a parsed icon and whether the set has it.
func Get(name string) (*ui.SVG, bool) {
	mu.Lock()
	defer mu.Unlock()
	if svg, ok := cache[name]; ok {
		return svg, svg != nil
	}
	data, err := read(name)
	if err != nil {
		cache[name] = nil
		return nil, false
	}
	svg, err := ui.ParseSVG(data)
	if err != nil {
		cache[name] = nil
		return nil, false
	}
	cache[name] = svg
	return svg, true
}

// Must returns a parsed icon and panics for a name the set lacks.
func Must(name string) *ui.SVG {
	svg, ok := Get(name)
	if !ok {
		panic("appicons: no icon named " + name)
	}
	return svg
}

// Names lists every icon in the set, sorted.
func Names() []string {
	seen := map[string]bool{}
	var names []string
	for _, dir := range dirs {
		entries, err := files.ReadDir(dir)
		if err != nil {
			panic("appicons: " + err.Error())
		}
		for _, e := range entries {
			name := strings.TrimSuffix(e.Name(), ".svg")
			if !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}
	sort.Strings(names)
	return names
}
