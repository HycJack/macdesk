package main

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func newTest(t *testing.T) (*desk, *ui.Tester) {
	t.Helper()
	isolateApps(t)
	d := newDesk()
	tt := ui.NewTester(d.view, 1280, 820)
	// Let arrival animations rest, or presses route nowhere.
	tt.Advance(300 * time.Millisecond)
	return d, tt
}

func TestDesktopLayout(t *testing.T) {
	_, tt := newTest(t)
	for _, s := range []string{"Macintosh HD", "会议记录"} {
		if !tt.HasText(s) {
			t.Errorf("missing %q; texts %v", s, tt.Texts())
		}
	}
	for _, s := range []string{"访达", "备忘录"} {
		if _, ok := tt.Find(s); !ok {
			t.Errorf("missing labeled element %q", s)
		}
	}
}

// TestFinderTilesStayInsideTheWindow keeps the icon grid from spilling out.
func TestFinderTilesStayInsideTheWindow(t *testing.T) {
	d, tt := newTest(t)
	d.wins["notes"].open = false
	d.wins["finder"].x, d.wins["finder"].y = 40, 64
	settle(tt)
	checkTiles := func(left, top, right, bottom float32, where string) {
		t.Helper()
		for _, f := range d.tr().finderFiles {
			r, ok := tt.Find(f.Name)
			if !ok {
				t.Fatalf("%s: tile %q not laid out; texts %v", where, f.Name, tt.Texts())
			}
			if r.X < left-0.5 || r.X+r.W > right+0.5 || r.Y < top-0.5 || r.Y+r.H > bottom+0.5 {
				t.Errorf("%s: tile %q at %+v leaves %v,%v-%v,%v", where, f.Name, r,
					left, top, right, bottom)
			}
		}
	}
	checkTiles(40, 64, 40+winFW, 64+winFH, "windowed")
	d.wins["finder"].zoom = true
	settle(tt)
	checkTiles(0, menuH, winW, winH-dockH, "zoomed")
}

func TestCloseAndReopenFinder(t *testing.T) {
	d, tt := newTest(t)
	if err := tt.Click("关闭"); err != nil {
		t.Fatal(err)
	}
	if d.wins["finder"].open {
		t.Fatal("finder still open")
	}
	if err := tt.Click("访达"); err != nil {
		t.Fatal(err)
	}
	if !d.wins["finder"].open {
		t.Error("dock did not reopen finder")
	}
}

func TestIconOpensFinder(t *testing.T) {
	d, tt := newTest(t)
	if err := tt.Click("关闭"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Macintosh HD"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Macintosh HD"); err != nil {
		t.Fatal(err)
	}
	if !d.wins["finder"].open {
		t.Error("double click did not open finder")
	}
}

func TestWindowDrag(t *testing.T) {
	d, tt := newTest(t)
	box, ok := tt.Find("项目")
	if !ok {
		t.Fatal("finder window not found")
	}
	// The labeled box is the window; press the top strip.
	cx, cy := box.X+box.W/2, box.Y+15
	before := d.wins["finder"].x
	tt.Press(cx, cy)
	tt.Move(cx+60, cy)
	tt.Release(cx+60, cy)
	if d.focus != "finder" {
		t.Errorf("drag did not focus the window: %q", d.focus)
	}
	if got := d.wins["finder"].x; got != before+60 {
		t.Errorf("window x %v, want %v", got, before+60)
	}
}

func TestLanguageToggle(t *testing.T) {
	d, tt := newTest(t)
	if err := tt.Click("EN"); err != nil {
		t.Fatal(err)
	}
	if d.lang != "en" || !tt.HasText("File") || tt.HasText("文件") {
		t.Errorf("lang %q texts %v", d.lang, tt.Texts())
	}
	if err := tt.Click("中"); err != nil {
		t.Fatal(err)
	}
	if d.lang != "zh" || !tt.HasText("文件") {
		t.Error("did not switch back")
	}
}

func TestMenuOpensFinder(t *testing.T) {
	d, tt := newTest(t)
	if err := tt.Click("关闭"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("文件"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("新建访达窗口"); err != nil {
		t.Fatal(err)
	}
	if !d.wins["finder"].open {
		t.Error("menu command did not open finder")
	}
}

func TestNotesSidebar(t *testing.T) {
	d, tt := newTest(t)
	if err := tt.Click("购物清单"); err != nil {
		t.Fatal(err)
	}
	if d.noteIn != "购物清单" {
		t.Errorf("note %q", d.noteIn)
	}
	// TextArea content is editor text, invisible to HasText.
	if d.noteText != zh.noteBodies[1] {
		t.Errorf("note body did not follow: %q", d.noteText)
	}
}

func TestDarkMode(t *testing.T) {
	d, tt := newTest(t)
	if !d.controls.Dark || !d.dark {
		t.Fatal("demo must start dark per spec")
	}
	if err := tt.Click("控制中心"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("深色模式"); err != nil {
		t.Fatal(err)
	}
	if d.controls.Dark || d.dark {
		t.Error("dark switch did not turn the theme light")
	}
	// The popover stays open; toggle back directly.
	if err := tt.Click("深色模式"); err != nil {
		t.Fatal(err)
	}
	if !d.controls.Dark || !d.dark {
		t.Error("dark switch did not restore the theme")
	}
}

func TestDefaultDark(t *testing.T) {
	d, _ := newTest(t)
	if !d.controls.Dark || !d.dark {
		t.Error("demo must start dark per spec")
	}
}

func TestAppleMenuHide(t *testing.T) {
	d, tt := newTest(t)
	box, ok := tt.Find("项目")
	if !ok {
		t.Fatal("finder window not found")
	}
	tt.Press(box.X+box.W/2, box.Y+15)
	tt.Release(box.X+box.W/2, box.Y+15)
	if d.focus != "finder" {
		t.Fatalf("focus %q", d.focus)
	}
	if err := tt.Click("Apple"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("隐藏"); err != nil {
		t.Fatal(err)
	}
	if d.wins["finder"].open {
		t.Error("hide did not close the focused window")
	}
}
