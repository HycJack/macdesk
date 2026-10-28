package desktop

import (
	"fmt"
	"math"
	"strings"

	"github.com/HycJack/macdesk/kit"
	"github.com/HycJack/macdesk/sound"
	"github.com/egoist/mygo/ui"
)

// DesktopIcon selects on click; opening bumps opened; cue Open.
func DesktopIcon(c *ui.Context, name, artName string, selected *bool, opened *int) ui.Element {
	if strings.TrimSpace(name) == "" || artName == "" || selected == nil || opened == nil || *opened < 0 {
		panic(fmt.Sprintf("desktop: DesktopIcon %q needs an icon, a selection and an open count of 0 or more", name))
	}
	t := c.Theme()
	b := ui.ButtonBase(c).Label(name).Width(t.Space(23)).Padding(t.Space(1)).Radius(t.Radius).AlignItems(ui.Center)
	kit.Arrive(b, ui.Pose{Scale: 0.85, Fade: 1}, 0)
	on := *selected
	b.Background(b.AnimateColor("bg", kit.Glass(t, 0.24*unit(b.Hovered() && !on)), ui.Snappy))
	b.Children(func() {
		ui.Column(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			ic := Art(c, artName, t.Space(16))
			ic.Radius(t.Radius).Border(2, ic.AnimateColor("ring", t.Accent.Alpha(unit(on)), ui.Snappy))
			label := ui.Text(c, name).FontSize(t.Rem(0.76)).TextAlign(ui.Center).MaxWidth(t.Space(20)).
				MaxLines(2).Padding(t.Space(0.25), t.Space(1.25)).Radius(t.Space(1))
			label.Background(label.AnimateColor("pill", t.Accent.Alpha(unit(on)), ui.Snappy)).
				TextColor(label.AnimateColor("ink", ui.Hex("#ffffff").Mix(t.AccentText, unit(on)), ui.Snappy))
		})
	})
	open := b.DoubleClicked() || (on && b.Shortcut(0, ui.KeyEnter))
	switch {
	case open && *opened < math.MaxInt:
		*selected = true
		*opened++
		b.SetChanged()
		sound.Play(sound.Open)
		kit.Log("DesktopIcon", "opened")
	case b.Clicked() && !on:
		*selected = true
		b.SetChanged()
		sound.Play(sound.Select)
		kit.Log("DesktopIcon", "selected")
	}
	return b
}
