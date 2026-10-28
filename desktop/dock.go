package desktop

import (
	"fmt"
	"strconv"
	"time"

	"github.com/HycJack/macdesk/appicons"
	"github.com/HycJack/macdesk/kit"
	"github.com/HycJack/macdesk/sound"
	"github.com/egoist/mygo/ui"
)

// DockApp is a dock app; Art names its icon and a positive Badge shows.
type DockApp struct {
	Name, Art string
	Running   bool
	Badge     int
	// Separator draws a divider before the icon, as before the trash.
	Separator bool
}

// dockSize is the icon square of the Dock, and dockPanel its height.
const (
	dockSize  float32 = 54
	dockPanel float32 = 74
)

// Dock magnifies under the pointer; clicks report; cue Tap.
func Dock(c *ui.Context, apps []DockApp, opened, right *string) ui.Element {
	names := make([]string, len(apps))
	for i, a := range apps {
		if a.Art == "" || a.Badge < 0 {
			panic(fmt.Sprintf("desktop: Dock app %q needs an icon and a badge of 0 or more", a.Name))
		}
		names[i] = a.Name
	}
	unique("Dock", names)
	if len(apps) == 0 || opened == nil {
		panic("desktop: Dock needs apps and an opened name")
	}
	t := c.Theme()
	root := ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.End).Padding(t.Space(1.75)).
		Radius(dockSize*0.4).Background(kit.Glass(t, 0.42)).Border(1, kit.Glass(t, 0.5)).
		Shadow(0, 18, 40, -8, kit.Shade(0.42)).Role(ui.RoleToolbar)
	kit.Arrive(root, ui.Pose{Y: t.Space(3), Fade: 1}, 0)
	hover := ui.Local(root, "hover", func() int { return -1 })
	root.Children(func() {
		sheen(c)
		now := -1
		for i, a := range apps {
			if a.Separator {
				ui.Box(c).Shrink(0).Width(1).Height(t.Space(9)).AlignSelf(ui.Center).
					Margin(t.Space(0.5), t.Space(1.5), t.Space(0.5), t.Space(1.5)).
					Background(kit.Glass(t, 0.34))
			}
			c.Key(strconv.Quote(a.Name))
			b := ui.ButtonBase(c).Label(a.Name).Column().Size(t.Space(14), dockPanel).
				AlignItems(ui.Center).Justify(ui.End).TransformOrigin(0.5, 1)
			if b.Hovered() {
				now = i
			}
			target := float32(1)
			if *hover >= 0 {
				switch d := max(i-*hover, *hover-i); d {
				case 0:
					target = 1.28
				case 1:
					target = 1.12
				}
			}
			scale := b.AnimateTo("scale", target, ui.Snappy)
			bounced := ui.Local(b, "bounced", func() time.Time { return time.Time{} })
			up := c.Now().Sub(*bounced) < 180*time.Millisecond
			if up {
				c.After(180*time.Millisecond - c.Now().Sub(*bounced))
			}
			lift := b.AnimateTo("lift", -t.Space(3.5)*unit(up), ui.Bouncy)
			b.Scale(scale, scale).Translate(0, lift)
			b.Children(func() {
				Art(c, a.Art, dockSize)
				if a.Badge > 0 {
					ui.Box(c).Absolute().Right(t.Space(1.5)).Top(t.Space(0.5)).Size(t.Space(5.5), t.Space(5.5)).
						Center().Radius(999).Background(t.Danger).
						Border(1.5, kit.Glass(t, 0.9)).Children(func() {
						ui.Text(c, strconv.Itoa(min(a.Badge, 99))).FontSize(t.Rem(0.58)).Bold().
							TextColor(t.AccentText).Shrink(0).SingleLine().ContentTransition(ui.ContentTransition{})
					})
				}
				ui.Box(c).Size(t.Space(1.25), t.Space(1.25)).Margin(t.Space(1), 0, t.Space(1.25), 0).
					Radius(999).Background(t.Text.Alpha(0.8 * unit(a.Running)))
			})
			if b.Clicked() {
				*opened = a.Name
				*bounced = c.Now()
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("Dock", "opened", "index", i)
			}
			if right != nil && b.RightClicked() {
				*right = a.Name
				root.SetChanged()
				kit.Log("Dock", "menu", "index", i)
			}
		}
		if now != *hover {
			*hover = now
			c.After(time.Millisecond)
		}
	})
	return root
}

// Art draws an app icon in the set, falling back to a plain document.
func Art(c *ui.Context, name string, size float32) ui.Element {
	svg, ok := appicons.Get(name)
	if !ok {
		svg = appicons.Must("f-app")
	}
	return ui.Image(c, svg).Size(size, size).Shadow(0, 2, 5, 0, kit.Shade(0.22))
}

// sheen is the specular light across the top of the glass.
func sheen(c *ui.Context) {
	ui.Box(c).Absolute().Left(0).Right(0).Top(0).Height(dockPanel * 0.46).
		LinearGradient(ui.LinearGradient{
			Angle: 180, From: ui.Hex("#ffffff").Alpha(0.26), To: ui.Hex("#ffffff").Alpha(0),
		})
}
