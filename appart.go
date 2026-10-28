package main

import "github.com/HycJack/macdesk/appicons"

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
