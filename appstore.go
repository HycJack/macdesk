package main

import (
	"maps"
	"slices"

	"github.com/HycJack/macdesk/desktop"
	"github.com/HycJack/macdesk/kit"
	"github.com/HycJack/macdesk/sound"
	"github.com/egoist/mygo/ui"
)

// storeBody is the App Store window: list, details and import.
func (d *desk) storeBody(c *ui.Context) {
	p := d.tr()
	t := c.Theme()
	ui.Row(c).Fill().Children(func() {
		ui.Box(c).Width(t.Space(40)).FillHeight().
			Background(t.Surface.Alpha(0.6)).Children(func() {
			sl := desktop.SourceList(c, d.storeGroups(p), &d.storeSel).FillHeight()
			if sl.Changed() {
				c.Invalidate()
			}
		})
		ui.Column(c).Grow(1).FillHeight().Padding(t.Space(3)).Gap(t.Space(2)).Children(func() {
			d.storeDetail(c, p)
			ui.Box(c).Grow(1)
			d.storeImport(c, p)
		})
	})
}

// storeGroups builds the installed and available sidebar groups.
func (d *desk) storeGroups(p *phrases) []desktop.SourceGroup {
	seen := map[string]bool{}
	var installed, available []desktop.SourceItem
	for _, id := range slices.Sorted(maps.Keys(d.installed)) {
		name := d.installed[id].manifest.Name
		if seen[name] {
			continue
		}
		seen[name] = true
		installed = append(installed, desktop.SourceItem{Name: name, Icon: "grid-2x2"})
	}
	for _, s := range samplePacks {
		if d.installed[s.id] != nil || seen[s.name] {
			continue
		}
		seen[s.name] = true
		available = append(available, desktop.SourceItem{Name: s.name, Icon: "download"})
	}
	var groups []desktop.SourceGroup
	if len(installed) > 0 {
		groups = append(groups, desktop.SourceGroup{Title: p.installedT, Items: installed})
	}
	if len(available) > 0 {
		groups = append(groups, desktop.SourceGroup{Title: p.availableT, Items: available})
	}
	if len(groups) == 0 {
		groups = append(groups, desktop.SourceGroup{Title: p.availableT,
			Items: []desktop.SourceItem{{Name: p.storeEmpty, Icon: "sparkles"}}})
	}
	var names []string
	for _, g := range groups {
		for _, it := range g.Items {
			names = append(names, it.Name)
		}
	}
	if !slices.Contains(names, d.storeSel) && len(names) > 0 {
		d.storeSel = names[0]
	}
	return groups
}

// storeDetail shows the selected app with its actions.
func (d *desk) storeDetail(c *ui.Context, p *phrases) {
	t := c.Theme()
	ia := d.installedByName(d.storeSel)
	sample := sampleByName(d.storeSel)
	if ia == nil && sample == nil {
		kit.Muted(c, p.storeEmpty).FillWidth().TextAlign(ui.Center)
		return
	}
	name, version, perms := "", "", []string{}
	if ia != nil {
		name, version, perms = ia.manifest.Name, ia.manifest.Version, ia.manifest.Permissions
	} else {
		name, version, perms = sample.name, sample.version, sample.perms
	}
	ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(t.Space(10), t.Space(10)).Center().Radius(t.Radius * 2).
				Background(t.Accent.Alpha(0.15)).Children(func() {
				glyph(c, "grid-2x2", t.Space(6)).TextColor(t.Accent)
			})
			ui.Column(c).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, name).FontSize(t.Rem(1.05)).Bold()
				kit.Muted(c, "v"+version).FontSize(t.Rem(0.72))
			})
		})
		if len(perms) > 0 {
			ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				kit.Muted(c, p.permsT).FontSize(t.Rem(0.72))
				for _, perm := range perms {
					ui.Text(c, perm).FontSize(t.Rem(0.72)).TextColor(t.Accent)
				}
			})
		}
		if d.storeErr != "" {
			ui.Text(c, p.storeErrT+"："+d.storeErr).FontSize(t.Rem(0.72)).
				TextColor(t.Danger).FillWidth().Wrap()
		}
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			d.storeButtons(c, p, ia, sample)
		})
	})
}

// storeButtons renders the open/install and uninstall buttons.
func (d *desk) storeButtons(c *ui.Context, p *phrases, ia *installedApp, s *samplePack) {
	main := p.openL
	if ia == nil {
		main = p.installL
	}
	b := ui.PrimaryButton(c, main)
	if kit.Clicked(kit.Press(b), sound.Tap) {
		if ia != nil {
			d.openID("app:" + ia.manifest.ID)
		} else if err := d.installSample(*s); err != nil {
			d.storeErr = err.Error()
		} else {
			d.storeErr = ""
		}
		c.Invalidate()
	}
	if ia != nil {
		if kit.Clicked(kit.Press(ui.Button(c, p.uninstallL)), sound.Tap) {
			if err := d.uninstallApp(ia.manifest.ID); err != nil {
				d.storeErr = err.Error()
			} else {
				d.storeErr = ""
			}
			c.Invalidate()
		}
	}
}

// storeImport installs a package from a typed path.
func (d *desk) storeImport(c *ui.Context, p *phrases) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
		ui.SearchField(c, &d.importIn).FillWidth().Placeholder(p.importPH)
		if kit.Clicked(kit.Press(ui.Button(c, p.importL)), sound.Tap) {
			if err := d.importApp(d.importIn); err != nil {
				d.storeErr = err.Error()
			} else {
				d.storeErr = ""
				d.importIn = ""
			}
			c.Invalidate()
		}
	})
}
