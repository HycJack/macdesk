package main

import (
	"strconv"

	"github.com/HycJack/macdesk/kit"
	"github.com/HycJack/macdesk/sound"
	"github.com/egoist/mygo/ui"
)

// settingsBody is the Settings window: appearance, wallpaper, accent, about.
func (d *desk) settingsBody(c *ui.Context) {
	p, t := d.tr(), c.Theme()
	ui.Column(c).Fill().Padding(t.Space(4)).Gap(t.Space(4)).Children(func() {
		d.setAppearance(c, p)
		d.setWallpapers(c, p)
		d.setAccents(c, p)
		card(c, p.setAbout, func() {
			kit.Muted(c, p.setVersion+" 1.0 · macdesk").FontSize(t.Rem(0.72))
		})
	})
}

// card wraps a settings section with its title.
func card(c *ui.Context, title string, fn func()) {
	t := c.Theme()
	ui.Column(c).FillWidth().Padding(t.Space(3)).Radius(t.Radius*1.5).Gap(t.Space(2)).
		Background(t.Surface).Border(1, hairline(t)).Children(func() {
		kit.Muted(c, title).FontSize(t.Rem(0.7)).Bold().SingleLine()
		fn()
	})
}

// setAppearance offers light and dark in a segmented pair.
func (d *desk) setAppearance(c *ui.Context, p *phrases) {
	t := c.Theme()
	card(c, p.setAppearance, func() {
		ui.Row(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(0.75)).
			Radius(t.Radius*1.5).Background(t.Background).Children(func() {
			for i, label := range [2]string{p.setLight, p.setDark} {
				on := (i == 1) == d.controls.Dark
				b := kit.Press(ui.ButtonBase(c).Label(label).Grow(1).Center().
					Padding(t.Space(1), 0).Radius(t.Radius))
				b.Background(b.AnimateColor("bg", t.Accent.Alpha(0.9*unit(on)), ui.Snappy))
				b.Children(func() {
					txt := ui.Text(c, label).FontSize(t.Rem(0.78)).SingleLine()
					if on {
						txt.TextColor(ui.Hex("#ffffff"))
					}
				})
				if kit.Clicked(b, sound.Tap) {
					d.controls.Dark = i == 1
					b.SetChanged()
					kit.Log("settings", "appearance", "dark", d.controls.Dark)
				}
			}
		})
	})
}

// setWallpapers shows one swatch per desktop picture.
func (d *desk) setWallpapers(c *ui.Context, p *phrases) {
	t := c.Theme()
	card(c, p.setWallpaper, func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			for i, pre := range wallpapers {
				i, pre := i, pre
				c.Key("wall:" + strconv.Itoa(i))
				b := kit.Press(ui.ButtonBase(c).Label(p.wallNames[i]).Column().
					Grow(1).Gap(t.Space(1)).Padding(t.Space(0.75)).Radius(t.Radius))
				b.Border(2, b.AnimateColor("ring", t.Accent.Alpha(unit(i == d.wall)), ui.Snappy))
				b.Children(func() {
					ui.Box(c).FillWidth().Height(t.Space(9)).Radius(t.Radius).
						LinearGradient(ui.LinearGradient{Angle: pre.angle, Stops: pre.stops, Oklab: true})
					kit.Muted(c, p.wallNames[i]).FontSize(t.Rem(0.64)).SingleLine()
				})
				if kit.Clicked(b, sound.Tap) {
					d.wall = i
					b.SetChanged()
					kit.Log("settings", "wallpaper", "index", i)
				}
			}
		})
	})
}

// setAccents shows the accent tones as dots.
func (d *desk) setAccents(c *ui.Context, p *phrases) {
	t := c.Theme()
	card(c, p.setAccent, func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			for i, tone := range accentTones {
				i, tone := i, tone
				c.Key("accent:" + strconv.Itoa(i))
				b := kit.Press(ui.ButtonBase(c).Label(p.accentNames[i]).Size(t.Space(7), t.Space(7)).Center())
				b.Radius(999).Background(tone.light).Scale(b.AnimateTo("s", 1+0.18*unit(i == d.accent), ui.Snappy),
					b.AnimateTo("s", 1+0.18*unit(i == d.accent), ui.Snappy))
				if i == d.accent {
					b.Border(2, t.Text.Alpha(0.85))
				}
				if kit.Clicked(b, sound.Tap) {
					d.accent = i
					b.SetChanged()
					kit.Log("settings", "accent", "index", i)
				}
			}
		})
	})
}
