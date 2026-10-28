package main

import (
	"testing"
)

func TestStoreInstallOpen(t *testing.T) {
	isolateApps(t)
	d, tt := newTest(t)
	if err := tt.Click("App Store"); err != nil {
		t.Fatal(err)
	}
	if !d.wins["store"].open {
		t.Fatal("store did not open")
	}
	if err := tt.Click("计算器"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("安装"); err != nil {
		t.Fatal(err)
	}
	if d.installed["calc"] == nil {
		t.Fatal("calc not installed")
	}
	if err := tt.Click("打开"); err != nil {
		t.Fatal(err)
	}
	if !d.wins["app:calc"].open {
		t.Error("calc window did not open")
	}
	if d.apps["calc"] == nil || d.apps["calc"].tree.Load() == nil {
		t.Error("calc instance not started")
	}
}

func TestStoreSelectionFallsBack(t *testing.T) {
	isolateApps(t)
	d, tt := newTest(t)
	if err := tt.Click("App Store"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("倒计时"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("安装"); err != nil {
		t.Fatal(err)
	}
	if d.installed["countdown"] == nil {
		t.Fatal("install failed")
	}
	// A dead selection falls back on the next frame.
	d.storeSel = "不存在"
	tt.Frame()
	if d.storeSel == "不存在" {
		t.Error("selection kept a dead item")
	}
}
