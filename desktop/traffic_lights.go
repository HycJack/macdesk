package desktop

import (
	"github.com/HycJack/macdesk/appicons"
	"github.com/HycJack/macdesk/kit"
	"github.com/HycJack/macdesk/sound"
	"github.com/egoist/mygo/ui"
)

// TrafficLightsOptions name the three buttons.
type TrafficLightsOptions struct{ Close, Minimize, Zoom string }

// TrafficLights close, minimize and zoom; cue Tap.
func TrafficLights(c *ui.Context, o TrafficLightsOptions) (root, close, minimize, zoom ui.Element) {
	def(&o.Close, "Close")
	def(&o.Minimize, "Minimize")
	def(&o.Zoom, "Zoom")
	t := c.Theme()
	root = ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleToolbar)
	kit.Arrive(root, ui.Pose{Scale: 0.8, Fade: 1}, 0)
	root.Children(func() {
		light := func(label, mark string, from, to ui.Color) ui.Element {
			b := kit.Press(ui.ButtonBase(c).Label(label).Size(12, 12).Center().Radius(999))
			b.LinearGradient(ui.LinearGradient{Angle: 180, From: from, To: to}).
				Border(0.5, kit.Shade(0.18))
			reveal := b.AnimateTo("reveal", unit(root.Hovered()), ui.Snappy)
			b.Children(func() {
				ui.Icon(c, appicons.Must(mark)).Size(t.Space(2.5), t.Space(2.5)).
					TextColor(kit.Shade(0.55)).Opacity(reveal)
			})
			if b.Clicked() {
				sound.Play(sound.Tap)
				kit.Log("TrafficLights", "pressed", "button", label)
			}
			return b
		}
		close = light(o.Close, "g-x", ui.Hex("#ff6a5f"), ui.Hex("#ee4a3f"))
		minimize = light(o.Minimize, "g-minus", ui.Hex("#ffd166"), ui.Hex("#f5b21c"))
		zoom = light(o.Zoom, "g-plus", ui.Hex("#5fdc7d"), ui.Hex("#28c840"))
	})
	return root, close, minimize, zoom
}
