package main

import (
	"slices"

	"github.com/HycJack/macdesk/appicons"
	"github.com/HycJack/macdesk/desktop"
	"github.com/HycJack/macdesk/kit"
	"github.com/HycJack/macdesk/sound"
	"github.com/egoist/mygo/ui"
)

// launchItem is one tile of Launchpad.
type launchItem struct{ name, art, win string }

// launchItems lists the built-in apps, then the installed packages.
func (d *desk) launchItems(p *phrases) []launchItem {
	items := []launchItem{
		{p.finder, "app-finder", "finder"},
		{p.safari, "app-safari", ""},
		{p.mail, "app-mail", ""},
		{p.notes, "app-notes", "notes"},
		{p.photos, "app-photos", ""},
		{p.music, "app-music", ""},
		{p.terminal, "app-terminal", ""},
		{p.settings, "app-settings", "settings"},
		{p.storeT, "app-appstore", "store"},
	}
	for _, id := range sortedApps(d.installed) {
		ia := d.installed[id]
		items = append(items, launchItem{ia.manifest.Name, appArt(id, ia.manifest.Icon), "app:" + id})
	}
	if q := lower(d.lpadQ); q != "" {
		items = slices.DeleteFunc(items, func(it launchItem) bool {
			return !containsFold(it.name, q)
		})
	}
	return items
}

// launchpad is the full screen grid of every app.
func (d *desk) launchpad(c *ui.Context, p *phrases) {
	t := c.Theme()
	layer := sheet(c).Background(kit.Shade(0.5))
	if layer.OverlayShortcut(0, ui.KeyEscape) {
		d.lpad = false
	}
	layer.Children(func() {
		scrim := sheet(c)
		if scrim.Clicked() {
			d.lpad = false
			sound.Play(sound.Tap)
		}
		ui.Column(c).Absolute().Top(t.Space(10)).Left(0).Right(0).
			AlignItems(ui.Center).Gap(t.Space(8)).Children(func() {
			d.launchSearch(c, p)
			d.launchGrid(c, p)
		})
	})
}

// launchSearch is the rounded search field above the grid.
func (d *desk) launchSearch(c *ui.Context, p *phrases) {
	t := c.Theme()
	field := ui.Row(c).Width(t.Space(105)).Height(t.Space(9)).Gap(t.Space(2)).
		Padding(0, t.Space(3)).Radius(t.Space(4.5)).
		AlignItems(ui.Center).Background(kit.Shade(0.22)).Border(1, kit.Glass(t, 0.28)).Shrink(0)
	field.Children(func() {
		ui.Icon(c, appicons.Must("g-search")).Size(t.Space(4), t.Space(4)).TextColor(ui.Hex("#ffffff").Alpha(0.8))
		in := ui.TextInputBase(c, &d.lpadQ).Placeholder(p.spotPH).FontSize(t.Rem(0.9)).
			TextColor(ui.Hex("#ffffff")).Grow(1).MinWidth(0)
		if in.Changed() {
			field.SetChanged()
		}
	})
}

// launchGrid lays the tiles out in rows of six.
func (d *desk) launchGrid(c *ui.Context, p *phrases) {
	t := c.Theme()
	items := d.launchItems(p)
	if len(items) == 0 {
		ui.Text(c, p.spotEmpty).FontSize(t.Rem(0.95)).TextColor(ui.Hex("#ffffff").Alpha(0.85))
		return
	}
	const perRow = 6
	ui.Column(c).Gap(t.Space(5)).Children(func() {
		for start := 0; start < len(items); start += perRow {
			row := items[start:min(start+perRow, len(items))]
			ui.Row(c).Gap(t.Space(2)).Justify(ui.Center).Children(func() {
				for _, it := range row {
					d.launchTile(c, it)
				}
			})
		}
	})
}

// launchTile is one app of the grid: artwork over a shadowed name.
func (d *desk) launchTile(c *ui.Context, it launchItem) {
	t := c.Theme()
	c.Key("lpad:" + it.name)
	tile := kit.Press(ui.ButtonBase(c).Label(it.name).Column().Width(t.Space(25)).Shrink(0).Clip().
		Gap(t.Space(1.5)).Padding(t.Space(2)).Radius(t.Space(3)).
		AlignItems(ui.Center).Justify(ui.Center))
	tile.Background(tile.AnimateColor("bg", kit.Shade(0.22*unit(tile.Hovered())), ui.Snappy))
	tile.Children(func() {
		desktop.Art(c, it.art, t.Space(16))
		ui.Text(c, it.name).FontSize(t.Rem(0.78)).TextColor(ui.Hex("#ffffff")).
			MaxLines(1).Ellipsis("").TextAlign(ui.Center).Shrink(0)
	})
	if kit.Clicked(tile, sound.Tap) {
		d.launchGo(it)
		tile.SetChanged()
	}
}

// launchGo closes the grid and opens the window of the chosen tile.
func (d *desk) launchGo(it launchItem) {
	d.lpad = false
	if it.win != "" {
		d.openID(it.win)
	}
	kit.Log("launchpad", "open", "app", it.name)
}
