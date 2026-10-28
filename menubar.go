package main

import (
	"strconv"

	"github.com/HycJack/macdesk/appicons"
	"github.com/HycJack/macdesk/desktop"
	"github.com/HycJack/macdesk/kit"
	"github.com/HycJack/macdesk/sound"
	"github.com/egoist/mygo/ui"
)

// menuBar is the top system bar with menus and status.
func (d *desk) menuBar(c *ui.Context, p *phrases) {
	t := c.Theme()
	bar := ui.Row(c).Absolute().Top(0).Left(0).Right(0).Height(menuH).AlignItems(ui.Center).
		Padding(0, t.Space(1)).Background(glass(t, 0.72)).
		BorderWidth(0, 0, 1, 0).BorderColor(hairline(t))
	bar.Children(func() {
		d.appleMenu(c, p)
		ui.Box(c).Grow(1).Children(func() {
			r := desktop.MenuBar(c, "macdesk", d.menus(p), &d.chosen, desktop.MenuBarOptions{Flat: true})
			if r.Changed() {
				d.menuChosen(p)
			}
		})
		desktop.WifiItem(c, d.networks(), &d.wifi,
			desktop.WifiItemOptions{Title: p.wifiTitle, Empty: p.wifiEmpty})
		desktop.BatteryItem(c, 0.82, desktop.BatteryItemOptions{Charging: true})
		desktop.ClockItem(c, desktop.ClockItemOptions{Format: p.clockFmt})
		d.statusButton(c, p, p.spotT, "g-spot", d.spot, func() { d.spot = !d.spot; d.spotQ = "" })
		d.statusButton(c, p, p.notifT, "g-cc", d.notifOpen, func() { d.notifOpen = !d.notifOpen })
		d.langButton(c, p)
		d.ccButton(c, p)
	})
}

// menus builds the app menu titles and commands.
func (d *desk) menus(p *phrases) []desktop.Menu {
	return []desktop.Menu{
		{Title: p.fileT, Items: []desktop.MenuItem{{Label: p.fileNew, Shortcut: "鈱楴"}}},
		{Title: p.editT, Items: []desktop.MenuItem{
			{Label: p.editUndo, Shortcut: "鈱榋"}, {Label: p.editRedo, Shortcut: "鈬р寴Z"},
		}},
		{Title: p.viewT, Items: []desktop.MenuItem{
			{Label: p.viewIcons, Shortcut: "鈱?"}, {Label: p.viewList, Shortcut: "鈱?"},
		}},
		{Title: p.windowT, Items: []desktop.MenuItem{
			{Label: p.winMin, Shortcut: "鈱楳"}, {Label: p.winZoom},
		}},
		{Title: p.helpT, Items: []desktop.MenuItem{{Label: p.helpOpen}}},
	}
}

// menuChosen runs the chosen menu command.
func (d *desk) menuChosen(p *phrases) {
	cmd := d.chosen
	if cmd == "" {
		return
	}
	d.chosen = ""
	switch cmd {
	case p.fileNew:
		d.openID("finder")
	case p.winMin:
		if ws := d.wins[d.focus]; ws != nil {
			ws.mini = true
		}
	case p.winZoom:
		if ws := d.wins[d.focus]; ws != nil {
			ws.zoom = !ws.zoom
		}
	default:
		kit.Log("macdesk", "menu", "cmd", cmd)
	}
}

// appleMenu is the system menu on the far left.
func (d *desk) appleMenu(c *ui.Context, p *phrases) {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label("Apple").
		Size(t.Space(5.5), t.Space(5.5)).Center().Radius(t.Radius))
	b.Children(func() { glyph(c, "applelogo", 16) })
	if kit.Clicked(b, sound.Tap) {
		d.apple = !d.apple
		b.SetChanged()
		if d.apple {
			kit.Log("macdesk", "apple")
		}
	}
	ui.PopoverBase(c, b, &d.apple, func(panel ui.Element) {
		panel.Margin(t.Space(1), 0, 0, 0).MinWidth(t.Space(52)).Padding(t.Space(1)).
			Gap(t.Space(0.5)).Background(t.Background).Border(1, t.Border).
			Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(0, 0)
		for _, it := range []desktop.MenuItem{
			{Label: p.appleAbout}, {Label: p.appleSettings}, {Label: p.appleHide}, {Label: p.appleQuit},
		} {
			c.Key(strconv.Quote(it.Label))
			item := kit.Press(ui.ButtonBase(c).Label(it.Label).Role(ui.RoleMenuItem).
				FillWidth().Padding(t.Space(1.25), t.Space(2)).Radius(t.Radius))
			item.Background(item.AnimateColor("bg", t.Accent.Alpha(unit(item.Hovered())), ui.Snappy))
			item.Children(func() {
				ui.Text(c, it.Label).FontSize(t.Rem(0.8)).SingleLine()
			})
			if kit.Clicked(item, sound.Tap) {
				d.apple = false
				switch it.Label {
				case p.appleAbout:
					d.openID("about")
				case p.appleHide, p.appleQuit:
					if ws := d.wins[d.focus]; ws != nil {
						ws.open = false
						d.focus = ""
					}
				default:
					kit.Log("macdesk", "apple", "cmd", it.Label)
				}
				item.SetChanged()
			}
		}
	})
}

// statusButton is a menu bar icon that toggles a panel.
func (d *desk) statusButton(c *ui.Context, p *phrases, label, art string, on bool, toggle func()) {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Size(t.Space(5.5), t.Space(5.5)).
		Center().Radius(t.Radius).Expanded(on))
	b.Children(func() { ui.Icon(c, appicons.Must(art)).Size(t.Space(4), t.Space(4)) })
	b.Background(b.AnimateColor("bg", t.Border.Alpha(0.6*max(unit(b.Hovered()), unit(on))), ui.Snappy))
	if kit.Clicked(b, sound.Tap) {
		toggle()
		b.SetChanged()
		kit.Log("macdesk", "status", "item", label)
	}
}

// langButton toggles Chinese and English copy.
func (d *desk) langButton(c *ui.Context, p *phrases) {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(p.langToggle).
		Padding(t.Space(1), t.Space(1.5)).Radius(t.Radius))
	b.Children(func() {
		ui.Text(c, p.langToggle).FontSize(t.Rem(0.78))
	})
	if kit.Clicked(b, sound.Tap) {
		d.toggleLang()
		b.SetChanged()
		kit.Log("macdesk", "lang", "to", d.lang)
	}
}

// networks is the simulated Wi-Fi environment.
func (d *desk) networks() []desktop.Network {
	return []desktop.Network{
		{Name: "Crychic-5G", Bars: 3, Secure: true},
		{Name: "Crychic-2.4G", Bars: 2, Secure: true},
		{Name: "Coffee-Guest", Bars: 1},
	}
}

// ccButton opens the Control Center popover.
func (d *desk) ccButton(c *ui.Context, p *phrases) {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(p.ccLabel).
		Size(t.Space(5.5), t.Space(5.5)).Center().Radius(t.Radius).Expanded(d.ccOpen))
	b.Children(func() { glyph(c, "sliders-horizontal", 16) })
	b.Background(b.AnimateColor("bg", t.Border.Alpha(0.6*max(unit(b.Hovered()), unit(d.ccOpen))), ui.Snappy))
	if kit.Clicked(b, sound.Tap) {
		d.ccOpen = !d.ccOpen
		b.SetChanged()
		if d.ccOpen {
			kit.Log("macdesk", "control-center")
		}
	}
	ui.PopoverBase(c, b, &d.ccOpen, func(panel ui.Element) {
		panel.Margin(t.Space(1), 0, 0, 0).Width(t.Space(78)).Padding(t.Space(2)).
			Background(t.Background).Border(1, t.Border).
			Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(1, 0)
		desktop.ControlCenter(c, &d.controls, desktop.ControlCenterOptions{
			Toggles: p.ccT, Volume: p.ccVol, Brightness: p.ccBri, On: p.ccOn, Off: p.ccOff,
		})
	})
}
