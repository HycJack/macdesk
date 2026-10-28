package main

import (
	"github.com/HycJack/macdesk/appicons"
	"github.com/HycJack/macdesk/desktop"
	"github.com/HycJack/macdesk/kit"
	"github.com/HycJack/macdesk/sound"
	"github.com/egoist/mygo/ui"
)

// spotPrefix keys the Spotlight rows for the tester.
const spotPrefix = "spot:"

// spotResult is one row of Spotlight's list.
type spotResult struct {
	group, name, kind, art, file string
	run                          func()
}

// spotResults matches the query against apps, actions and documents.
func (d *desk) spotResults(p *phrases) []spotResult {
	q := lower(d.spotQ)
	var apps, acts, files []spotResult
	for _, it := range d.launchItems(p) {
		win := it.win
		apps = append(apps, spotResult{group: p.spotApps, name: it.name, kind: p.spotApps, art: it.art,
			run: func() {
				if win != "" {
					d.openID(win)
				}
			}})
	}
	for _, a := range []struct {
		name, art string
		run       func()
	}{
		{p.actDark, "g-moon", func() { d.controls.Dark = !d.controls.Dark }},
		{p.actWall, "g-globe", func() { d.wall = (d.wall + 1) % len(wallpapers) }},
		{p.actLpad, "app-launchpad", func() { d.lpad = true }},
		{p.actNotif, "g-cc", func() { d.notifOpen = true }},
		{p.actTrash, "app-trash", d.clearTrash},
		{p.actAbout, "g-info", func() { d.openID("about") }},
	} {
		acts = append(acts, spotResult{group: p.spotActions, name: a.name, kind: p.spotActions, art: a.art,
			run: a.run})
	}
	for _, f := range p.finderFiles {
		files = append(files, spotResult{group: p.spotFiles, name: f.Name, kind: f.Meta, art: f.Mark,
			file: f.Name, run: func() { d.openID("finder") }})
	}
	keep := func(in []spotResult) []spotResult {
		if q == "" {
			return in
		}
		var out []spotResult
		for _, it := range in {
			if containsFold(it.name, q) {
				out = append(out, it)
			}
		}
		return out
	}
	return append(append(keep(apps), keep(acts)...), keep(files)...)
}

// spotlight is the search overlay in the middle of the screen.
func (d *desk) spotlight(c *ui.Context, p *phrases) {
	t := c.Theme()
	layer := sheet(c)
	if layer.OverlayShortcut(0, ui.KeyEscape) {
		d.spot = false
	}
	if layer.OverlayShortcut(0, ui.KeyDown) {
		d.spotSel++
	}
	if layer.OverlayShortcut(0, ui.KeyUp) {
		d.spotSel = max(0, d.spotSel-1)
	}
	layer.Children(func() {
		scrim := sheet(c).Background(kit.Shade(0.26))
		if scrim.Clicked() {
			d.spot = false
			sound.Play(sound.Tap)
		}
		panel := ui.Column(c).Absolute().Top(menuH + t.Space(22)).Left(t.Space(40)).Right(t.Space(40)).
			MaxWidth(t.Space(160)).AlignItems(ui.Center)
		panel.Children(func() {
			box := ui.Column(c).Width(t.Space(148)).Background(glass(t, 0.92)).
				Radius(t.Space(3.5)).Border(1, hairline(t)).
				Shadow(0, 24, 60, -6, kit.Shade(0.4))
			box.Children(func() {
				d.spotField(c, p, box)
				d.spotList(c, p)
			})
		})
	})
}

// spotField is the rounded query row.
func (d *desk) spotField(c *ui.Context, p *phrases, box ui.Element) {
	t := c.Theme()
	ui.Row(c).FillWidth().Height(t.Space(11)).Gap(t.Space(2)).AlignItems(ui.Center).
		Padding(0, t.Space(3.5)).Children(func() {
		ui.Icon(c, appicons.Must("g-spot")).Size(t.Space(5), t.Space(5)).TextColor(t.TextMuted)
		in := ui.TextInputBase(c, &d.spotQ).Placeholder(p.spotPH).FontSize(t.Rem(1.15)).
			Grow(1).MinWidth(0).AutoFocus()
		if in.Changed() {
			d.spotPick(0)
			box.SetChanged()
		}
		if in.Submitted() {
			d.spotFire(p)
		}
	})
	ui.Box(c).FillWidth().Height(1).Background(hairline(t))
}

// spotList draws the grouped results.
func (d *desk) spotList(c *ui.Context, p *phrases) {
	t := c.Theme()
	res := d.spotResults(p)
	if len(res) == 0 {
		ui.Row(c).FillWidth().Padding(t.Space(3)).Justify(ui.Center).Children(func() {
			kit.Muted(c, p.spotEmpty).FontSize(t.Rem(0.85))
		})
		return
	}
	if d.spotSel >= len(res) {
		d.spotSel = 0
	}
	group := ""
	ui.Column(c).FillWidth().Padding(t.Space(1.5)).Gap(t.Space(0.25)).Children(func() {
		for i, r := range res {
			if r.group != group {
				group = r.group
				kit.Muted(c, group).FontSize(t.Rem(0.66)).Padding(t.Space(1), t.Space(2))
			}
			d.spotRow(c, p, r, i)
		}
	})
}

// spotRow is one clickable result.
func (d *desk) spotRow(c *ui.Context, p *phrases, r spotResult, i int) {
	t := c.Theme()
	c.Key(spotPrefix + r.name)
	row := kit.Press(ui.ButtonBase(c).Label(r.name).Row().FillWidth().Height(t.Space(7.5)).
		Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(0, t.Space(2)).Radius(t.Radius))
	on := i == d.spotSel
	row.Background(kit.Shade(0.10 * max(unit(row.Hovered()), unit(on))))
	if on {
		row.Background(t.Accent.Alpha(0.9))
	}
	row.Children(func() {
		desktop.Art(c, r.art, t.Space(5.5))
		ui.Text(c, r.name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
		if on {
			ui.Text(c, r.kind).FontSize(t.Rem(0.7)).TextColor(ui.Hex("#ffffff").Alpha(0.8)).SingleLine()
		} else {
			kit.Muted(c, r.kind).FontSize(t.Rem(0.7)).SingleLine()
		}
	})
	if row.Hovered() {
		d.spotPick(i)
	}
	if kit.Clicked(row, sound.Tap) {
		d.spotPick(i)
		d.spotFire(p)
	}
}

// spotPick selects a result and keeps the list in step.
func (d *desk) spotPick(i int) {
	if d.spotSel == i {
		return
	}
	d.spotSel = i
}

// spotFire runs the selected result.
func (d *desk) spotFire(p *phrases) {
	res := d.spotResults(p)
	if len(res) == 0 {
		return
	}
	sel := clampInt(d.spotSel, 0, len(res)-1)
	d.spot = false
	d.spotQ = ""
	res[sel].run()
	kit.Log("spotlight", "run", "name", res[sel].name)
}
