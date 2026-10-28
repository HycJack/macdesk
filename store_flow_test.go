package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStoreFullFlow(t *testing.T) {
	d, tt := newTest(t)
	// Open the store from the dock.
	if err := tt.Click("App Store"); err != nil {
		t.Fatal(err)
	}
	// Install the snake.
	if err := tt.Click("贪吃蛇"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("安装"); err != nil {
		t.Fatal(err)
	}
	if d.installed["snake"] == nil {
		t.Fatal("snake not installed")
	}
	// Open it and play one move.
	if err := tt.Click("打开"); err != nil {
		t.Fatal(err)
	}
	if d.apps["snake"] == nil {
		t.Fatal("snake instance missing")
	}
	if err := tt.Click("→"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(2 * time.Second)
	if inst := d.apps["snake"]; inst.lastErr != "" {
		t.Fatalf("snake errored: %s", inst.lastErr)
	}
	// Close the window; the instance stops.
	d.wins["app:snake"].open = false
	d.stopApp("snake")
	tt.Frame()
	// Uninstall from the store.
	if err := tt.Click("贪吃蛇"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("卸载"); err != nil {
		t.Fatal(err)
	}
	if d.installed["snake"] != nil {
		t.Fatal("snake still installed")
	}
	if _, err := os.Stat(filepath.Join(appsRoot, "snake")); !os.IsNotExist(err) {
		t.Error("snake directory not removed")
	}
}
