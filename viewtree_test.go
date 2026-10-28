package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func openDemoApp(t *testing.T, d *desk, tt *ui.Tester) {
	t.Helper()
	if err := d.installApp(writePackage(t, t.TempDir(), "demo")); err != nil {
		t.Fatal(err)
	}
	d.wins["app:demo"] = &winState{open: true, x: 100, y: 100}
	tt.Frame()
}

func TestRenderAppView(t *testing.T) {
	d, tt := newTest(t)
	openDemoApp(t, d, tt)
	if !tt.HasText("n=0") {
		t.Errorf("app view not rendered: %v", tt.Texts())
	}
	d.apps["demo"].doAction(map[string]any{"type": "add"})
	tt.Frame()
	if !tt.HasText("n=1") {
		t.Error("tree did not refresh after action")
	}
}

func TestRenderLimits(t *testing.T) {
	d, tt := newTest(t)
	openDemoApp(t, d, tt)
	inst := d.apps["demo"]
	kids := make([]*Node, 501)
	for i := range kids {
		kids[i] = Text("x", "body")
	}
	inst.tree.Store(Column(0, kids...))
	tt.Frame()
	if !tt.HasText("应用异常") {
		t.Error("over-limit tree rendered anyway")
	}
}

func TestRenderButtonAction(t *testing.T) {
	d, tt := newTest(t)
	openDemoApp(t, d, tt)
	// Drive Action through the host path.
	d.appAction(d.apps["demo"], map[string]any{"type": "add"})
	tt.Frame()
	if !tt.HasText("n=1") {
		t.Error("host action did not refresh the view")
	}
}
