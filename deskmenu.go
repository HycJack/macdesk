package main

import (
	"github.com/HycJack/macdesk/appicons"
	"github.com/HycJack/macdesk/kit"
	"github.com/HycJack/macdesk/sound"
	"github.com/egoist/mygo/ui"
)

// deskMenu is the desktop's own context menu at the pointer.
func (d *desk) deskMenu(c *ui.Context, p *phrases) {
	t := c.Theme()
	scrim := sheet(c)
	if scrim.Clicked() || scrim.RightClicked() {
		d.ctxOpen = false
	}
	if scrim.OverlayShortcut(0, ui.KeyEscape) {
		d.ctxOpen = false
	}
	scrim.Children(func() {
		panel := ui.Column(c).Absolute().Left(d.ctxX).Top(d.ctxY).Width(t.Space(46)).
			Padding(t.Space(1)).Gap(t.Space(0.25)).Radius(t.Radius*1.5).
			Background(glass(t, 0.94)).Border(1, hairline(t)).
			Shadow(0, 14, 34, -4, kit.Shade(0.36))
		panel.Children(func() {
			for _, it := range []struct{ label, art string }{
				{p.ctxNewFolder, "g-docs"},
				{p.ctxWall, "g-globe"},
				{p.ctxLpad, "app-launchpad"},
				{p.ctxSpot, "g-spot"},
				{p.ctxTrash, "app-trash"},
				{p.ctxAbout, "g-info"},
			} {
				d.deskMenuItem(c, p, it.label, it.art)
			}
		})
	})
}

// deskMenuItem is one row of the desktop menu.
func (d *desk) deskMenuItem(c *ui.Context, p *phrases, label, art string) {
	t := c.Theme()
	row := kit.Press(ui.ButtonBase(c).Label(label).Row().FillWidth().Gap(t.Space(2)).
		AlignItems(ui.Center).Padding(t.Space(1), t.Space(1.5)).Radius(t.Radius))
	row.Background(row.AnimateColor("bg", t.Accent.Alpha(unit(row.Hovered())), ui.Snappy))
	row.Children(func() {
		ui.Image(c, appicons.Must(art)).Size(t.Space(3.5), t.Space(3.5)).Shrink(0)
		ui.Text(c, label).FontSize(t.Rem(0.78)).SingleLine()
	})
	if kit.Clicked(row, sound.Tap) {
		d.ctxOpen = false
		switch label {
		case p.ctxWall:
			d.wall = (d.wall + 1) % len(wallpapers)
		case p.ctxLpad:
			d.lpad = true
		case p.ctxSpot:
			d.spot = true
		case p.ctxTrash:
			d.clearTrash()
		case p.ctxAbout:
			d.openID("about")
		case p.ctxNewFolder:
			d.openID("finder")
		}
		row.SetChanged()
		kit.Log("desktop", "menu", "cmd", label)
	}
}

// winID maps a dock app name to its window id, empty when it has none.
func (d *desk) winID(p *phrases, name string) string {
	switch name {
	case p.finder:
		return "finder"
	case p.notes:
		return "notes"
	case p.storeT:
		return "store"
	case p.settings:
		return "settings"
	case p.trashT:
		return "trash"
	}
	for id, ia := range d.installed {
		if ia.manifest.Name == name {
			return "app:" + id
		}
	}
	return ""
}

// openDockMenu anchors the icon menu above the pointer.
func (d *desk) openDockMenu(c *ui.Context, p *phrases, at ui.Element) {
	if x, y, ok := at.PointerPosition(); ok {
		d.ctxX, d.ctxY = x-90, y-160
	}
	d.ctxOpen = true
	d.ctxApp = d.dockRight
	d.dockRight = ""
	kit.Log("Dock", "context", "app", d.ctxApp)
}

// dockMenu is the icon menu: open, close, quit.
func (d *desk) dockMenu(c *ui.Context, p *phrases) {
	t := c.Theme()
	scrim := sheet(c)
	if scrim.Clicked() || scrim.RightClicked() {
		d.ctxOpen = false
		d.ctxApp = ""
	}
	if scrim.OverlayShortcut(0, ui.KeyEscape) {
		d.ctxOpen = false
		d.ctxApp = ""
	}
	id := d.winID(p, d.ctxApp)
	scrim.Children(func() {
		panel := ui.Column(c).Absolute().Left(max(0, d.ctxX)).Top(max(0, d.ctxY)).Width(t.Space(44)).
			Padding(t.Space(1)).Gap(t.Space(0.25)).Radius(t.Radius*1.5).
			Background(glass(t, 0.94)).Border(1, hairline(t)).
			Shadow(0, 14, 34, -4, kit.Shade(0.36))
		panel.Children(func() {
			d.menuRow(c, p.dockOpenL, func() {
				if id != "" {
					d.openID(id)
				}
			})
			if id != "" {
				d.menuRow(c, p.dockCloseL, func() {
					if ws := d.wins[id]; ws != nil {
						ws.open = false
					}
				})
			}
			if len(id) > 4 && id[:4] == "app:" {
				d.menuRow(c, p.dockQuitL, func() {
					if ws := d.wins[id]; ws != nil {
						ws.open = false
					}
					d.stopApp(id[4:])
				})
			}
		})
	})
}

// menuRow is one row of a popup menu.
func (d *desk) menuRow(c *ui.Context, label string, run func()) {
	t := c.Theme()
	row := kit.Press(ui.ButtonBase(c).Label(label).Row().FillWidth().Gap(t.Space(2)).
		AlignItems(ui.Center).Padding(t.Space(1), t.Space(1.5)).Radius(t.Radius))
	row.Background(row.AnimateColor("bg", t.Accent.Alpha(unit(row.Hovered())), ui.Snappy))
	row.Children(func() { ui.Text(c, label).FontSize(t.Rem(0.78)).SingleLine() })
	if kit.Clicked(row, sound.Tap) {
		d.ctxOpen = false
		d.ctxApp = ""
		run()
		row.SetChanged()
		kit.Log("macdesk", "menu", "cmd", label)
	}
}
