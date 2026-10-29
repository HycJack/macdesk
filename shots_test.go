package main

import (
	"image/png"
	"os"
	"testing"
	"time"

	"github.com/HycJack/macdesk/appicons"
	"github.com/egoist/mygo/ui"
)

// TestShootDocs renders the README screenshots into docs/.
// Run it with MACDESK_SHOTS=1 go test -run TestShootDocs .
func TestShootDocs(t *testing.T) {
	if os.Getenv("MACDESK_SHOTS") == "" {
		t.Skip("set MACDESK_SHOTS=1 to refresh docs screenshots")
	}
	isolateApps(t)
	d := newDesk()
	tt := ui.NewTester(d.view, 1280, 820)
	settle(tt)
	shoot(t, tt, "desktop")
	if err := d.installSample(*sampleByID("countdown")); err != nil {
		t.Fatal(err)
	}
	d.wins["finder"].open = false
	d.wins["notes"].open = false
	d.openID("store")
	settle(tt)
	shoot(t, tt, "store")
	d.wins["store"].open = false
	d.openID("app:countdown")
	settle(tt)
	shoot(t, tt, "app")
}

// TestShootOverlays renders Launchpad, Spotlight, the tray and Settings.
func TestShootOverlays(t *testing.T) {
	if os.Getenv("MACDESK_SHOTS") == "" {
		t.Skip("set MACDESK_SHOTS=1 to refresh docs screenshots")
	}
	isolateApps(t)
	d := newDesk()
	for _, id := range []string{"calc", "snake"} {
		if err := d.installSample(*sampleByID(id)); err != nil {
			t.Fatal(err)
		}
	}
	d.wins["notes"].open = false
	tt := ui.NewTester(d.view, 1280, 820)
	settle(tt)
	d.lpad = true
	settle(tt)
	shoot(t, tt, "launchpad")
	d.lpad = false
	d.spot = true
	d.spotQ = "设"
	settle(tt)
	shoot(t, tt, "spotlight")
	d.spot = false
	d.notifOpen = true
	settle(tt)
	shoot(t, tt, "notifications")
	d.notifOpen = false
	d.openID("settings")
	settle(tt)
	shoot(t, tt, "settings")
}

// TestShootFinder renders the Finder icon view for the README.
func TestShootFinder(t *testing.T) {
	if os.Getenv("MACDESK_SHOTS") == "" {
		t.Skip("set MACDESK_SHOTS=1 to refresh docs screenshots")
	}
	isolateApps(t)
	d := newDesk()
	d.wins["notes"].open = false
	tt := ui.NewTester(d.view, 1280, 820)
	settle(tt)
	d.wins["finder"].x, d.wins["finder"].y = 40, 64
	settle(tt)
	shoot(t, tt, "finder")
	d.wins["finder"].zoom = true
	settle(tt)
	shoot(t, tt, "finder-wide")
}

// TestShootApps renders the Photos and Terminal windows for the README.
func TestShootApps(t *testing.T) {
	if os.Getenv("MACDESK_SHOTS") == "" {
		t.Skip("set MACDESK_SHOTS=1 to refresh docs screenshots")
	}
	isolateApps(t)
	d := newDesk()
	d.wins["notes"].open = false
	tt := ui.NewTester(d.view, 1280, 820)
	settle(tt)
	d.openID("photos")
	d.openID("terminal")
	d.wins["photos"].x, d.wins["photos"].y = 60, 70
	d.wins["terminal"].x, d.wins["terminal"].y = 600, 300
	settle(tt)
	shoot(t, tt, "apps")
}

// TestShootLight renders the light appearance for the README.
func TestShootLight(t *testing.T) {
	if os.Getenv("MACDESK_SHOTS") == "" {
		t.Skip("set MACDESK_SHOTS=1 to refresh docs screenshots")
	}
	isolateApps(t)
	d := newDesk()
	d.controls.Dark = false
	if err := d.installSample(*sampleByID("calc")); err != nil {
		t.Fatal(err)
	}
	d.wins["notes"].open = false
	d.openID("store")
	tt := ui.NewTester(d.view, 1280, 820)
	settle(tt)
	shoot(t, tt, "light")
}

// TestShootIconSheet renders the icon set for the README.
func TestShootIconSheet(t *testing.T) {
	if os.Getenv("MACDESK_SHOTS") == "" {
		t.Skip("set MACDESK_SHOTS=1 to refresh docs screenshots")
	}
	names := appicons.Names()
	const cols, cell = 13, 56
	rows := (len(names) + cols - 1) / cols
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Gap(0).Background(ui.Hex("#eceef2")).Children(func() {
			for r := 0; r < rows; r++ {
				ui.Row(c).Gap(0).Children(func() {
					for i := r * cols; i < min((r+1)*cols, len(names)); i++ {
						ui.Box(c).Size(cell, cell).Center().Children(func() {
							ui.Image(c, appicons.Must(names[i])).Size(44, 44)
						})
					}
				})
			}
		})
	}, cols*cell, rows*cell)
	tt.Frame()
	f, err := os.Create("docs/icons.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, tt.Image()); err != nil {
		t.Fatal(err)
	}
}

// shoot writes the tester's frame to docs/name.png.
func shoot(t *testing.T, tt *ui.Tester, name string) {
	t.Helper()
	tt.Frame()
	f, err := os.Create("docs/" + name + ".png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, tt.Image()); err != nil {
		t.Fatal(err)
	}
}

// settle runs a few frames so arrivals paint.
func settle(tt *ui.Tester) {
	for range 8 {
		tt.Frame()
		tt.Advance(120 * time.Millisecond)
	}
}
