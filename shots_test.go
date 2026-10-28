package main

import (
	"image/png"
	"os"
	"testing"
	"time"

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
