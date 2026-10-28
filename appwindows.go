package main

import (
	"github.com/HycJack/macdesk/appicons"
	"github.com/HycJack/macdesk/desktop"
	"github.com/HycJack/macdesk/kit"
	"github.com/HycJack/macdesk/sound"
	"github.com/egoist/mygo/ui"
)

// mailRows are the sample messages of the Mail window.
var mailRows = []struct{ from, subject, when string }{
	{"Apple", "您的 Apple ID 已在新设备登录", "09:24"},
	{"GitHub", "macdesk: 3 new stars", "08:51"},
	{"林可", "周会纪要已更新", "昨天"},
	{"App Store", "贪吃蛇 已更新到 1.0.1", "昨天"},
}

// photoTiles label the grid of the Photos window.
var photoTiles = []string{"日落", "海岸", "街头", "雪山", "夜市", "森林",
	"城市", "花火", "天空", "咖啡", "书店", "雨天"}

// termLines are the fixed transcript of the Terminal window.
var termLines = []string{
	"macdesk 1.0.0 · go 1.27 · windows/amd64",
	"$ ls apps",
	"calc  countdown  snake  sudoku",
	"$ goapp run calc",
	"loading app… ok",
	"$ ",
}

// safariBody is a start page: a toolbar and a row of favourites.
func (d *desk) safariBody(c *ui.Context) {
	p, t := d.tr(), c.Theme()
	ui.Column(c).Fill().Children(func() {
		desktop.UnifiedToolbar(c, p.safariStart, &d.safariIn, desktop.UnifiedToolbarOptions{
			Back: p.backL, Forward: p.fwdL, Placeholder: p.safariPH,
		})
		ui.Column(c).Grow(1).FillWidth().Padding(t.Space(5)).Gap(t.Space(3)).Children(func() {
			kit.Muted(c, p.favT).FontSize(t.Rem(0.72)).Bold().SingleLine()
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				for _, f := range []struct{ name, art string }{
					{p.safari, "app-safari"}, {p.mail, "app-mail"}, {p.storeT, "app-appstore"},
					{p.notes, "app-notes"}, {p.photos, "app-photos"}, {p.finder, "app-finder"},
				} {
					card := kit.Press(ui.ButtonBase(c).Label(f.name).Column().Grow(1).
						Gap(t.Space(1.5)).Padding(t.Space(2.5)).Radius(t.Radius * 1.5))
					card.Background(card.AnimateColor("bg", t.Surface, ui.Snappy)).
						Border(1, hairline(t))
					card.Children(func() {
						ui.Image(c, appicons.Must(f.art)).Size(t.Space(8), t.Space(8))
						ui.Text(c, f.name).FontSize(t.Rem(0.72)).MaxLines(1).
							Ellipsis("").TextAlign(ui.Center).SingleLine()
					})
					if kit.Clicked(card, sound.Tap) {
						d.openApp(p, f.name)
					}
				}
			})
		})
	})
}

// mailBody is the mailbox sidebar and the message list.
func (d *desk) mailBody(c *ui.Context) {
	p, t := d.tr(), c.Theme()
	groups := []desktop.SourceGroup{
		{Title: p.mailBoxes, Items: []desktop.SourceItem{
			{Name: p.mailInbox, Icon: "g-download", Count: len(mailRows)},
			{Name: p.mailSent, Icon: "g-fwd", Count: 0},
			{Name: p.mailDrafts, Icon: "g-compose", Count: 0},
		}},
	}
	ui.Row(c).Fill().Children(func() {
		ui.Column(c).Width(t.Space(50)).FillHeight().Padding(t.Space(2)).Gap(t.Space(2)).
			Background(glass(t, 0.5)).BorderWidth(0, 0, 0, 1).BorderColor(hairline(t)).
			Children(func() {
				ui.Text(c, p.mail).FontSize(t.Rem(0.82)).Bold().Padding(t.Space(1), t.Space(2))
				ui.Box(c).Grow(1).Children(func() {
					if desktop.SourceList(c, groups, &d.mailBox).Changed() {
						c.Invalidate()
					}
				})
			})
		ui.Column(c).Grow(1).FillHeight().Children(func() {
			for _, m := range mailRows {
				row := kit.Press(ui.ButtonBase(c).Label(m.subject).Row().FillWidth().
					Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)))
				row.Background(row.AnimateColor("bg",
					t.Surface.Alpha(max(0.55, unit(row.Hovered()))), ui.Snappy))
				row.Children(func() {
					ui.Box(c).Size(t.Space(6), t.Space(6)).Radius(999).Shrink(0).
						Background(t.Accent.Alpha(0.18)).Center().Children(func() {
						ui.Text(c, m.from[:1]).FontSize(t.Rem(0.72)).Bold().TextColor(t.Accent)
					})
					ui.Column(c).Grow(1).MinWidth(0).Gap(1).Children(func() {
						ui.Text(c, m.from).FontSize(t.Rem(0.76)).Bold().SingleLine()
						kit.Muted(c, m.subject).FontSize(t.Rem(0.72)).MaxLines(1).Ellipsis("")
					})
					kit.Muted(c, m.when).FontSize(t.Rem(0.66)).Shrink(0)
				})
				ui.Box(c).FillWidth().Height(1).Background(hairline(t))
			}
		})
	})
}

// photosBody is a grid of tinted tiles.
func (d *desk) photosBody(c *ui.Context) {
	p, t := d.tr(), c.Theme()
	ui.Column(c).Fill().Padding(t.Space(3)).Gap(t.Space(2)).Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			ui.Text(c, p.photos).FontSize(t.Rem(0.85)).Bold().Grow(1).SingleLine()
			kit.Muted(c, "12").FontSize(t.Rem(0.7))
		})
		ui.Column(c).Grow(1).FillWidth().Gap(t.Space(2)).Children(func() {
			for start := 0; start < len(photoTiles); start += 4 {
				ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
					for i := start; i < min(start+4, len(photoTiles)); i++ {
						pre := wallpapers[i%len(wallpapers)]
						ui.Box(c).Grow(1).Height(t.Space(21)).Radius(t.Radius).Clip().
							LinearGradient(ui.LinearGradient{Angle: pre.angle, Stops: pre.stops, Oklab: true}).
							Children(func() {
								ui.Box(c).Absolute().Left(t.Space(1.5)).Bottom(t.Space(1.5)).
									Padding(2, t.Space(1)).Radius(t.Radius).Background(kit.Shade(0.32)).
									Children(func() {
										ui.Text(c, photoTiles[i]).FontSize(t.Rem(0.64)).
											TextColor(ui.Hex("#ffffff")).SingleLine()
									})
							})
					}
				})
			}
		})
	})
}

// musicBody is the player card.
func (d *desk) musicBody(c *ui.Context) {
	p, t := d.tr(), c.Theme()
	ui.Column(c).Fill().Center().Gap(t.Space(4)).Padding(t.Space(5)).Children(func() {
		ui.Box(c).Size(t.Space(46), t.Space(46)).Radius(t.Space(4)).Clip().
			LinearGradient(ui.LinearGradient{Angle: 150, Stops: wallpapers[2].stops, Oklab: true}).
			Shadow(0, 18, 40, -8, kit.Shade(0.4))
		ui.Column(c).Gap(t.Space(0.5)).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "Ocean Drift").FontSize(t.Rem(1.1)).Bold().SingleLine()
			kit.Muted(c, p.musicNow+" · macdesk sessions").FontSize(t.Rem(0.75)).SingleLine()
		})
		ui.Row(c).Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			for _, b := range []struct{ art, label string }{
				{"g-shuffle", p.musicNow}, {"g-bwd", p.backL}, {"g-play", p.playL},
				{"g-fwd", p.fwdL}, {"g-repeat", p.musicNow},
			} {
				btn := kit.Press(ui.ButtonBase(c).Label(b.label).Size(t.Space(7), t.Space(7)).Center())
				btn.Radius(999)
				if b.art == "g-play" {
					btn.Background(t.Accent)
				}
				btn.Children(func() {
					ui.Icon(c, appicons.Must(b.art)).Size(t.Space(4), t.Space(4))
				})
				if kit.Clicked(btn, sound.Tap) && b.art == "g-play" {
					kit.Log("music", "play", "track", "Ocean Drift")
				}
			}
		})
		ui.Box(c).Width(t.Space(46)).Height(t.Space(1)).Radius(999).
			Background(t.Text.Alpha(0.12)).Children(func() {
			ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Width(t.Space(20)).
				Background(t.Accent)
		})
	})
}

// terminalBody is the fixed transcript a shell would print.
func (d *desk) terminalBody(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Fill().Padding(t.Space(3)).Gap(t.Space(0.5)).
		Background(ui.Hex("#101014")).Children(func() {
		for i, line := range termLines {
			txt := ui.Text(c, line).FontSize(t.Rem(0.78)).SingleLine().
				TextColor(ui.Hex("#e6e6e8"))
			if i == 0 {
				txt.TextColor(ui.Hex("#8f8f96"))
			}
			if line == "$ " {
				txt.TextColor(ui.Hex("#59d98b"))
			}
		}
	})
}
