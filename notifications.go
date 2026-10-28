package main

import (
	"time"

	"github.com/HycJack/macdesk/appicons"
	"github.com/HycJack/macdesk/kit"
	"github.com/HycJack/macdesk/sound"
	"github.com/egoist/mygo/ui"
)

// notif is one entry of the Notification Centre.
type notif struct{ title, body, at string }

// pushNotif prepends a notification and keeps the list short.
func (d *desk) pushNotif(title, body string) {
	d.notifs = append([]notif{{title: title, body: body, at: time.Now().Format("15:04")}}, d.notifs...)
	if len(d.notifs) > 6 {
		d.notifs = d.notifs[:6]
	}
}

// notifTray is the panel that slides in from the right edge.
func (d *desk) notifTray(c *ui.Context, p *phrases) {
	t := c.Theme()
	layer := sheet(c)
	if layer.OverlayShortcut(0, ui.KeyEscape) {
		d.notifOpen = false
	}
	layer.Children(func() {
		scrim := sheet(c)
		if scrim.Clicked() {
			d.notifOpen = false
			sound.Play(sound.Tap)
		}
		panel := ui.Column(c).Absolute().Right(t.Space(2)).Top(menuH + t.Space(2)).
			Width(t.Space(82)).MaxHeight(winH - menuH - dockH - t.Space(6)).Gap(t.Space(2.5))
		kit.Arrive(panel, ui.Pose{X: t.Space(6), Fade: 1}, 0)
		panel.Children(func() {
			d.clockCard(c, p)
			d.notifCard(c, p)
		})
	})
}

// clockCard is the big time and the week strip.
func (d *desk) clockCard(c *ui.Context, p *phrases) {
	t := c.Theme()
	now := time.Now()
	ui.Column(c).FillWidth().Padding(t.Space(3)).Radius(t.Space(4)).
		Background(glass(t, 0.62)).Border(1, hairline(t)).Gap(t.Space(1.5)).Children(func() {
		ui.Text(c, now.Format("15:04")).FontSize(t.Rem(2.6)).Bold().SingleLine()
		kit.Muted(c, now.Format("2006-01-02")+" 路 "+p.today).FontSize(t.Rem(0.75)).SingleLine()
		ui.Row(c).FillWidth().Gap(t.Space(1)).Children(func() {
			for i, w := range p.calWeek {
				on := int(now.Weekday()+6)%7 == i
				cell := ui.Column(c).Grow(1).Center().Padding(t.Space(0.75), 0).Radius(t.Space(2)).Gap(2)
				if on {
					cell.Background(t.Accent)
				}
				cell.Children(func() {
					kit.Muted(c, w).FontSize(t.Rem(0.6))
					ui.Text(c, now.AddDate(0, 0, i-int(now.Weekday()+6)%7).Format("02")).
						FontSize(t.Rem(0.7)).SingleLine()
				})
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, appicons.Must("g-cloud")).Size(t.Space(5), t.Space(5)).TextColor(t.Accent)
			kit.Muted(c, p.weatherCity+" 路 "+p.weatherState).FontSize(t.Rem(0.72)).SingleLine()
		})
	})
}

// notifCard lists the notifications with a clear button.
func (d *desk) notifCard(c *ui.Context, p *phrases) {
	t := c.Theme()
	ui.Column(c).FillWidth().Padding(t.Space(3)).Radius(t.Space(4)).Gap(t.Space(2)).
		Background(glass(t, 0.62)).Border(1, hairline(t)).Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			ui.Text(c, p.notifT).FontSize(t.Rem(0.8)).Bold().Grow(1).SingleLine()
			if len(d.notifs) > 0 {
				b := kit.Press(ui.ButtonBase(c).Label(p.notifClear).Padding(t.Space(0.5), t.Space(1)).
					Radius(t.Radius))
				b.Children(func() { kit.Muted(c, p.notifClear).FontSize(t.Rem(0.7)) })
				if kit.Clicked(b, sound.Tap) {
					d.notifs = nil
					b.SetChanged()
				}
			}
		})
		if len(d.notifs) == 0 {
			kit.Muted(c, p.notifNone).FontSize(t.Rem(0.74))
			return
		}
		for _, n := range d.notifs[:min(len(d.notifs), 4)] {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start).Children(func() {
				ui.Icon(c, appicons.Must("app-appstore")).Size(t.Space(4), t.Space(4)).Shrink(0)
				ui.Column(c).Grow(1).MinWidth(0).Gap(1).Children(func() {
					ui.Text(c, n.title).FontSize(t.Rem(0.74)).Bold().SingleLine()
					kit.Muted(c, n.body).FontSize(t.Rem(0.68)).MaxLines(2).Wrap()
				})
				kit.Muted(c, n.at).FontSize(t.Rem(0.62)).Shrink(0)
			})
		}
	})
}
