package main

import (
	"maps"
	"slices"
	"strconv"

	"github.com/HycJack/macdesk/desktop"
	"github.com/HycJack/macdesk/kit"
	"github.com/HycJack/macdesk/sound"
	"github.com/egoist/mygo/ui"
)

// storeBody is the App Store window: sidebar, the chosen app, its facts and
// the import strip.
func (d *desk) storeBody(c *ui.Context) {
	p, t := d.tr(), c.Theme()
	ui.Row(c).Fill().Children(func() {
		ui.Column(c).Width(t.Space(52)).FillHeight().Padding(t.Space(2)).Gap(t.Space(2)).
			Background(glass(t, 0.5)).BorderWidth(0, 0, 0, 1).BorderColor(hairline(t)).
			Children(func() {
				ui.Text(c, p.storeT).FontSize(t.Rem(0.82)).Bold().Padding(t.Space(1), t.Space(2))
				ui.Box(c).Grow(1).Children(func() {
					if desktop.SourceList(c, d.storeGroups(p), &d.storeSel).Changed() {
						c.Invalidate()
					}
				})
			})
		ui.Column(c).Grow(1).FillHeight().Padding(t.Space(4)).Gap(t.Space(3)).Children(func() {
			d.storeHero(c, p)
			d.storeFacts(c, p)
			d.storeShelf(c, p)
			ui.Box(c).Grow(1)
			d.storeImport(c, p)
		})
	})
}

// storePick is the app the store shows, installed or offered.
type storePick struct {
	name, version, id, art, entry string
	api                           int
	perms                         []string
}

// pick gathers what the store knows about the chosen app.
func (d *desk) pick() *storePick {
	if ia := d.installedByName(d.storeSel); ia != nil {
		m := ia.manifest
		return &storePick{name: m.Name, version: m.Version, id: m.ID, art: appArt(m.ID, m.Icon),
			entry: m.entry(), api: m.API, perms: m.Permissions}
	}
	if s := sampleByName(d.storeSel); s != nil {
		return &storePick{name: s.name, version: s.version, id: s.id, art: appArt(s.id, ""),
			entry: "main.go", api: 1, perms: s.perms}
	}
	return nil
}

// storeGroups builds the installed and available sidebar groups.
func (d *desk) storeGroups(p *phrases) []desktop.SourceGroup {
	seen := map[string]bool{}
	var installed, available []desktop.SourceItem
	for _, id := range slices.Sorted(maps.Keys(d.installed)) {
		m := d.installed[id].manifest
		if seen[m.Name] {
			continue
		}
		seen[m.Name] = true
		installed = append(installed, desktop.SourceItem{Name: m.Name, Icon: appArt(m.ID, m.Icon)})
	}
	for _, s := range samplePacks {
		if d.installed[s.id] != nil || seen[s.name] {
			continue
		}
		seen[s.name] = true
		available = append(available, desktop.SourceItem{Name: s.name, Icon: appArt(s.id, "")})
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
			Items: []desktop.SourceItem{{Name: p.storeEmpty, Icon: "g-apps"}}})
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

// storeHero is the app's banner: artwork, name, version and its actions.
func (d *desk) storeHero(c *ui.Context, p *phrases) {
	t := c.Theme()
	pick := d.pick()
	if pick == nil {
		kit.Muted(c, p.storeEmpty).FillWidth().TextAlign(ui.Center)
		return
	}
	ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Center).Children(func() {
		desktop.Art(c, pick.art, t.Space(17))
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.75)).Children(func() {
			ui.Text(c, pick.name).FontSize(t.Rem(1.6)).Bold().SingleLine()
			kit.Muted(c, "v"+pick.version+" · "+pick.id).FontSize(t.Rem(0.8)).SingleLine()
		})
		d.storeButtons(c, p, d.installedByName(d.storeSel), sampleByName(d.storeSel))
	})
}

// storeFacts is the tagline, the package facts and the permissions.
func (d *desk) storeFacts(c *ui.Context, p *phrases) {
	t := c.Theme()
	pick := d.pick()
	if pick == nil {
		return
	}
	ui.Box(c).FillWidth().Padding(t.Space(3)).Radius(t.Radius*1.5).
		Background(t.Surface).Border(1, hairline(t)).Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(2.5)).Children(func() {
			kit.Muted(c, p.storeTag).FontSize(t.Rem(0.8)).FillWidth().Wrap()
			ui.Row(c).FillWidth().Gap(t.Space(6)).Children(func() {
				fact(c, p.storeEntryT, pick.entry)
				fact(c, p.storeAPIT, strconv.Itoa(pick.api))
				fact(c, p.storeLimitT, "2 MiB")
			})
			if len(pick.perms) > 0 {
				ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					kit.Muted(c, p.permsT).FontSize(t.Rem(0.74))
					for _, perm := range pick.perms {
						ui.Text(c, perm).FontSize(t.Rem(0.72)).TextColor(t.Accent).
							Padding(t.Space(0.5), t.Space(1.5)).Radius(999).
							Background(t.Accent.Alpha(0.12))
					}
				})
			}
			if d.storeErr != "" {
				ui.Text(c, p.storeErrT+"："+d.storeErr).FontSize(t.Rem(0.74)).
					TextColor(t.Danger).FillWidth().Wrap()
			}
		})
	})
}

// fact is one labelled package detail.
func fact(c *ui.Context, label, value string) {
	t := c.Theme()
	ui.Column(c).Gap(t.Space(0.25)).Children(func() {
		kit.Muted(c, label).FontSize(t.Rem(0.66)).SingleLine()
		ui.Text(c, value).FontSize(t.Rem(0.78)).SingleLine()
	})
}

// storeShelf is the row of other apps the store offers; a card picks one.
func (d *desk) storeShelf(c *ui.Context, p *phrases) {
	t := c.Theme()
	var offers []*samplePack
	for i := range samplePacks {
		s := &samplePacks[i]
		if d.installed[s.id] == nil && s.name != d.storeSel {
			offers = append(offers, s)
		}
	}
	if len(offers) == 0 {
		return
	}
	kit.Muted(c, p.moreT).FontSize(t.Rem(0.72)).Bold().SingleLine()
	ui.Row(c).FillWidth().Gap(t.Space(2.5)).Children(func() {
		for _, s := range offers[:min(len(offers), 4)] {
			c.Key("shelf:" + s.id)
			b := kit.Press(ui.ButtonBase(c).Label(s.name).Row().Grow(1).Gap(t.Space(2)).
				AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius*1.5).
				Background(t.Surface).Border(1, hairline(t)))
			b.Children(func() {
				desktop.Art(c, appArt(s.id, ""), t.Space(8))
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
					ui.Text(c, s.name).FontSize(t.Rem(0.8)).Bold().SingleLine()
					kit.Muted(c, "v"+s.version).FontSize(t.Rem(0.68)).SingleLine()
				})
			})
			if b.Clicked() {
				d.storeSel = s.name
				c.Invalidate()
				sound.Play(sound.Select)
				kit.Log("store", "shelf", "pick", s.id)
			}
		}
	})
}

// storeButtons renders the open/install and uninstall buttons.
func (d *desk) storeButtons(c *ui.Context, p *phrases, ia *installedApp, s *samplePack) {
	t := c.Theme()
	main := p.openL
	if ia == nil {
		main = p.installL
	}
	ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(t.Space(1.5)).Children(func() {
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
	})
}

// storeImport installs a package from a typed path.
func (d *desk) storeImport(c *ui.Context, p *phrases) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
		ui.TextInput(c, &d.importIn).FillWidth().Placeholder(p.importPH)
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
