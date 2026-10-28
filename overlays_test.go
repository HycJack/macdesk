package main

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// TestLaunchpadListsEveryApp checks the grid and its search filter.
func TestLaunchpadListsEveryApp(t *testing.T) {
	d, _ := newTest(t)
	p := d.tr()
	items := d.launchItems(p)
	names := make([]string, len(items))
	for i, it := range items {
		names[i] = it.name
	}
	joined := strings.Join(names, ",")
	for _, want := range []string{p.finder, p.settings, p.storeT, p.music} {
		if !strings.Contains(joined, want) {
			t.Errorf("launchpad lacks %q; has %v", want, names)
		}
	}
	if strings.Contains(joined, p.trashT) {
		t.Error("the trash belongs to the dock, not the grid")
	}
	d.lpadQ = "设置"
	got := d.launchItems(p)
	if len(got) != 1 || got[0].name != p.settings {
		t.Fatalf("search for 设置 gave %v", got)
	}
	d.lpadQ = "没有这个应用"
	if n := len(d.launchItems(p)); n != 0 {
		t.Fatalf("a bogus query kept %d items", n)
	}
}

// TestLaunchpadOpensSettings fills the grid and checks a tile opens a window.
func TestLaunchpadOpensSettings(t *testing.T) {
	d, tt := newTest(t)
	d.lpad = true
	settle(tt)
	if !tt.HasText(d.tr().settings) {
		t.Fatalf("launchpad did not render; texts %v", tt.Texts())
	}
	if !tt.HasText(d.tr().trashT) {
		t.Error("the grid misses the trash tile")
	}
	d.launchGo(launchItem{name: d.tr().settings, win: "settings"})
	settle(tt)
	if d.lpad {
		t.Error("launchpad stayed open after a tile was chosen")
	}
	if !d.wins["settings"].open {
		t.Error("settings window did not open")
	}
}

// TestSpotlightGroupsResults checks apps, actions and files group apart.
func TestSpotlightGroupsResults(t *testing.T) {
	d, _ := newTest(t)
	p := d.tr()
	d.spotQ = "废纸篓"
	res := d.spotResults(p)
	if len(res) == 0 {
		t.Fatal("no results for a known action")
	}
	if res[0].group != p.spotActions && res[0].group != p.spotApps {
		t.Errorf("unexpected group %q", res[0].group)
	}
	d.spotQ = "会议"
	res = d.spotResults(p)
	if len(res) != 1 || res[0].name != "会议记录.md" {
		t.Fatalf("document search gave %v", res)
	}
	d.spotQ = ""
	if all := d.spotResults(p); len(all) < 10 {
		t.Errorf("an empty query lists only %d results", len(all))
	}
}

// TestSpotlightFiresAction runs the selected result.
func TestSpotlightFiresAction(t *testing.T) {
	d, tt := newTest(t)
	d.spot = true
	d.spotQ = d.tr().actDark
	tt.Advance(200 * time.Millisecond)
	d.spotFire(d.tr())
	if d.spot {
		t.Error("spotlight stayed open")
	}
	if d.spotQ != "" {
		t.Errorf("query kept as %q", d.spotQ)
	}
	if d.controls.Dark {
		t.Error("the dark mode action did not flip the theme")
	}
}

// TestNotificationCentreKeepsInstalls records what the store does.
func TestNotificationCentreKeepsInstalls(t *testing.T) {
	d, tt := newTest(t)
	if err := d.installSample(*sampleByID("calc")); err != nil {
		t.Fatal(err)
	}
	if len(d.notifs) != 1 {
		t.Fatalf("notifications %v", d.notifs)
	}
	if !strings.Contains(d.notifs[0].title, d.tr().notifInstall) {
		t.Errorf("title %q lacks the install word", d.notifs[0].title)
	}
	d.notifOpen = true
	settle(tt)
	if !tt.HasText(d.tr().notifT) {
		t.Errorf("tray did not render; texts %v", tt.Texts())
	}
	if !tt.HasText(d.tr().notifClear) {
		t.Error("no clear button in the tray")
	}
}

// TestTrashKeepsUninstalledApps checks the bin and its empty action.
func TestTrashKeepsUninstalledApps(t *testing.T) {
	d, tt := newTest(t)
	if err := d.installSample(*sampleByID("snake")); err != nil {
		t.Fatal(err)
	}
	if err := d.uninstallApp("snake"); err != nil {
		t.Fatal(err)
	}
	if len(d.trash) != 1 || d.trash[0].name != "贪吃蛇" {
		t.Fatalf("trash %v", d.trash)
	}
	if trashArt(d.trash) != "app-trash-full" {
		t.Error("the dock still shows the empty bin")
	}
	d.openID("trash")
	settle(tt)
	if !tt.HasText("贪吃蛇") {
		t.Errorf("trash window misses the app; texts %v", tt.Texts())
	}
	d.clearTrash()
	if len(d.trash) != 0 {
		t.Error("clearTrash left items")
	}
	if trashArt(d.trash) != "app-trash" {
		t.Error("the dock kept the full bin")
	}
}

// TestSettingsSwitchAppearanceAndAccent drives the Settings window.
func TestSettingsSwitchAppearanceAndAccent(t *testing.T) {
	d, tt := newTest(t)
	d.openID("settings")
	settle(tt)
	if !tt.HasText(d.tr().setWallpaper) {
		t.Fatalf("settings did not render; texts %v", tt.Texts())
	}
	if err := tt.Click(d.tr().setLight); err != nil {
		t.Fatal(err)
	}
	settle(tt)
	if d.controls.Dark {
		t.Error("light was clicked but the desk stayed dark")
	}
	d.wall = 2
	d.accent = 4
	settle(tt)
	if accent(d.accent, false) != accentTones[4].light {
		t.Error("accent lookup drifted")
	}
}

// TestDesktopContextMenu opens on a secondary click.
func TestDesktopContextMenu(t *testing.T) {
	d, tt := newTest(t)
	tt.RightClickAt(300, 620)
	settle(tt)
	if !d.ctxOpen {
		t.Fatal("the desktop menu did not open")
	}
	if !tt.HasText(d.tr().ctxLpad) {
		t.Errorf("menu items missing; texts %v", tt.Texts())
	}
	if d.ctxX < 200 || d.ctxX > 400 || d.ctxY < 540 || d.ctxY > 700 {
		t.Errorf("menu opened at %v,%v instead of the pointer", d.ctxX, d.ctxY)
	}
	tt.Key(0, ui.KeyEscape)
	settle(tt)
	if d.ctxOpen {
		t.Error("escape did not close the menu")
	}
}

// TestDockMenuOpensAndCloses drives the dock icon menu.
func TestDockMenuOpensAndCloses(t *testing.T) {
	d, tt := newTest(t)
	d.wins["finder"].open = false
	settle(tt)
	rect, ok := tt.Find(d.tr().finder)
	if !ok {
		t.Fatal("no finder icon in the dock")
	}
	tt.RightClickAt(rect.X+rect.W/2, rect.Y+rect.H/2)
	settle(tt)
	if d.ctxApp != d.tr().finder {
		t.Fatalf("menu for %q, want the finder", d.ctxApp)
	}
	if !tt.HasText(d.tr().dockOpenL) {
		t.Fatalf("menu items missing; texts %v", tt.Texts())
	}
	if err := tt.Click(d.tr().dockOpenL); err != nil {
		t.Fatal(err)
	}
	settle(tt)
	if !d.wins["finder"].open {
		t.Error("the dock menu did not reopen the window")
	}
	if d.ctxOpen {
		t.Error("the menu stayed open")
	}
}

// TestDockOpensEveryBuiltInWindow walks the decorative apps.
func TestDockOpensEveryBuiltInWindow(t *testing.T) {
	d, tt := newTest(t)
	p := d.tr()
	for _, c := range []struct{ name, marker string }{
		{p.safari, p.favT},
		{p.mail, p.mailSent},
		{p.photos, "12"},
		{p.music, p.musicNow},
		{p.terminal, "macdesk 1.0.0 · go 1.27"},
	} {
		d.openApp(p, c.name)
		settle(tt)
		if !tt.HasText(c.marker) {
			t.Errorf("%s window misses %q; texts %v", c.name, c.marker, tt.Texts())
		}
	}
}

// TestSystemShortcuts binds Spotlight and Launchpad to keys.
func TestSystemShortcuts(t *testing.T) {
	d, tt := newTest(t)
	tt.Key(ui.Super, ui.KeySpace)
	settle(tt)
	if !d.spot {
		t.Fatal("super+space did not open Spotlight")
	}
	tt.Key(0, ui.KeyEscape)
	settle(tt)
	if d.spot {
		t.Error("escape did not close Spotlight")
	}
	tt.Key(0, ui.KeyF4)
	settle(tt)
	if !d.lpad {
		t.Fatal("F4 did not open Launchpad")
	}
	tt.Key(0, ui.KeyF4)
	settle(tt)
	if d.lpad {
		t.Error("F4 did not close Launchpad again")
	}
}
