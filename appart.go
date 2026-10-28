package main

import (
	"maps"
	"slices"
	"strings"

	"github.com/HycJack/macdesk/appicons"
)

// sampleArt maps the bundled apps to the icons they wear.
var sampleArt = map[string]string{
	"calc":      "app-calculator",
	"countdown": "app-clock",
	"snake":     "app-snake",
	"sudoku":    "app-sudoku",
}

// appArt names the icon of an app: the one its manifest names when the set
// has it, else the bundled one, else a plain app.
func appArt(id, named string) string {
	if named != "" {
		if _, ok := appicons.Get(named); ok {
			return named
		}
	}
	if art, ok := sampleArt[id]; ok {
		return art
	}
	return "f-app"
}

// sortedApps lists installed app ids in order.
func sortedApps(byID map[string]*installedApp) []string {
	return slices.Sorted(maps.Keys(byID))
}

// trashArt is the full bin when the trash holds something.
func trashArt(items []trashItem) string {
	if len(items) > 0 {
		return "app-trash-full"
	}
	return "app-trash"
}

// lower lowercases for a search box.
func lower(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

// containsFold reports whether hay holds needle, ignoring case.
func containsFold(hay, needle string) bool {
	return strings.Contains(strings.ToLower(hay), needle)
}
