package main

import (
	"time"

	"github.com/HycJack/macdesk/appicons"
	"github.com/HycJack/macdesk/kit"
	"github.com/HycJack/macdesk/sound"
	"github.com/egoist/mygo/ui"
)

// trashItem is an app that was uninstalled.
type trashItem struct{ name, at string }

// pushTrash remembers an uninstalled app.
func (d *desk) pushTrash(name string) {
	d.trash = append([]trashItem{{name: name, at: time.Now().Format("15:04")}}, d.trash...)
}

// clearTrash empties the bin.
func (d *desk) clearTrash() {
	if len(d.trash) == 0 {
		return
	}
	d.trash = nil
	sound.Play(sound.Tap)
	kit.Log("trash", "emptied")
}

// trashBody is the Trash window.
func (d *desk) trashBody(c *ui.Context) {
	p, t := d.tr(), c.Theme()
	ui.Column(c).Fill().Padding(t.Space(4)).Gap(t.Space(3)).Center().Children(func() {
		if len(d.trash) == 0 {
			ui.Image(c, appicons.Must("app-trash")).Size(t.Space(18), t.Space(18))
			ui.Text(c, p.trashEmptyT).FontSize(t.Rem(1.05)).Bold().SingleLine()
			kit.Muted(c, p.trashEmptyHint).FontSize(t.Rem(0.76)).SingleLine()
			return
		}
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			ui.Text(c, p.trashHint).FontSize(t.Rem(0.8)).Bold().Grow(1).SingleLine()
			if kit.Clicked(kit.Press(ui.Button(c, p.trashClear)), sound.Tap) {
				d.clearTrash()
				c.Invalidate()
			}
		})
		for _, it := range d.trash {
			ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).
				Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).
				Background(t.Surface).Children(func() {
				ui.Image(c, appicons.Must("f-app")).Size(t.Space(5), t.Space(5)).Shrink(0)
				ui.Text(c, it.name).FontSize(t.Rem(0.8)).Grow(1).SingleLine()
				kit.Muted(c, it.at).FontSize(t.Rem(0.66)).Shrink(0)
			})
		}
	})
}
