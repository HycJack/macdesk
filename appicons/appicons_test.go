package appicons

import "testing"

// TestEveryIconParses keeps the set loadable: one bad file breaks a window.
func TestEveryIconParses(t *testing.T) {
	names := Names()
	if len(names) < 100 {
		t.Fatalf("the set has only %d icons", len(names))
	}
	for _, name := range names {
		if _, ok := Get(name); !ok {
			t.Errorf("%s does not parse", name)
		}
	}
}

// TestAppIconsAreArtwork pins the icons the dock and the store wear.
func TestAppIconsAreArtwork(t *testing.T) {
	for _, name := range []string{"app-finder", "app-appstore", "app-calculator", "app-snake", "app-sudoku", "f-app"} {
		if _, ok := Get(name); !ok {
			t.Errorf("missing %s", name)
		}
	}
}

// TestUnknownIcon checks the fallback path of a name the set lacks.
func TestUnknownIcon(t *testing.T) {
	if _, ok := Get("no-such-icon"); ok {
		t.Error("an unknown name is reported as found")
	}
	defer func() {
		if recover() == nil {
			t.Error("Must did not panic for an unknown name")
		}
	}()
	Must("no-such-icon")
}
