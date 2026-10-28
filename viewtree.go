package main

import (
	"strconv"
	"strings"

	"github.com/HycJack/macdesk/icon"
	"github.com/HycJack/macdesk/kit"
	"github.com/HycJack/macdesk/sound"
	"github.com/egoist/mygo/ui"
)

// renderAppView draws an app's tree, its error or a placeholder.
func (d *desk) renderAppView(c *ui.Context, inst *appInstance) {
	p := d.tr()
	t := c.Theme()
	inst.maybeTick(c)
	if inst.broken != "" {
		kit.Muted(c, inst.broken).FillWidth().TextAlign(ui.Center)
		return
	}
	if inst.lastErr != "" {
		ui.Column(c).Fill().Center().Gap(t.Space(1)).Children(func() {
			ui.Text(c, p.appErrT).FontSize(t.Rem(0.9)).Bold()
			ui.Text(c, inst.lastErr).FontSize(t.Rem(0.72)).TextColor(t.TextMuted).
				FillWidth().Wrap().MaxWidth(t.Space(60))
		})
		return
	}
	root := inst.tree.Load()
	if root == nil {
		kit.Muted(c, "…").FillWidth().TextAlign(ui.Center)
		return
	}
	if !treeOK(root) {
		kit.Muted(c, p.appErrT).FillWidth().TextAlign(ui.Center)
		return
	}
	ui.Column(c).FillWidth().FillHeight().Padding(appPad).Gap(0).Children(func() {
		d.renderNode(c, inst, root, 0, "0")
	})
}

// appPad is the inset around an app's view tree, in DIPs.
const appPad = 14

// renderNode renders one node and its children.
func (d *desk) renderNode(c *ui.Context, inst *appInstance, n *Node, depth int, path string) {
	if depth > 12 {
		return
	}
	t := c.Theme()
	switch n.Type {
	case "column":
		ui.Column(c).FillWidth().Gap(t.Space(float32(n.Gap))).Children(func() {
			d.renderChildren(c, inst, n, depth, path)
		})
	case "row":
		ui.Row(c).FillWidth().Gap(t.Space(float32(n.Gap))).AlignItems(ui.Center).
			Children(func() {
				d.renderChildren(c, inst, n, depth, path)
			})
	case "box":
		box := ui.Box(c).Radius(t.Radius)
		if len(n.Children) == 0 {
			box.Size(t.Space(3), t.Space(3))
		} else {
			box.FillWidth().Padding(t.Space(2))
		}
		box.Background(appColor(t, n.Color))
		if len(n.Children) > 0 {
			box.Children(func() { d.renderChildren(c, inst, n, depth, path) })
		}
	case "text":
		size, bold := textSize(t, n.Size)
		txt := ui.Text(c, n.Text).FontSize(size).FillWidth().Wrap().MinWidth(0)
		if bold {
			txt.Bold()
		}
	case "button":
		b := kit.Press(ui.ButtonBase(c).Label(n.Label).
			Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius))
		if n.Primary {
			b.Background(t.Accent)
		}
		b.Children(func() {
			ui.Text(c, n.Label).FontSize(t.Rem(0.85)).SingleLine()
		})
		if kit.Clicked(b, sound.Tap) {
			d.sendNodeAction(inst, n.Action)
			c.Invalidate()
		}
	case "input":
		v := inst.field(path, "input", n.Text).Interface().(*string)
		in := ui.TextInput(c, v).FillWidth()
		if n.Label != "" {
			in.Placeholder(n.Label)
		}
		if in.Changed() {
			d.appAction(inst, map[string]any{"text": *v})
			c.Invalidate()
		}
	case "image":
		if knownIcons[n.Name] {
			glyph(c, n.Name, float32(n.Value)).TextColor(t.Text)
		} else {
			kit.Muted(c, "?")
		}
	case "progress":
		frac := float32(0)
		if n.Max > 0 {
			frac = float32(n.Value / n.Max)
		}
		ui.Column(c).FillWidth().Gap(0).Children(func() {
			ui.Box(c).FillWidth().Height(t.Space(1.5)).Radius(999).
				Background(t.Border).Children(func() {
				ui.Box(c).Width(t.Space(30 * frac)).Height(t.Space(1.5)).
					Radius(999).Background(t.Accent)
			})
		})
	case "slider":
		v := inst.field(path, "slider", n.Value).Interface().(*float64)
		s := ui.Slider(c, v, n.Min, n.Max).FillWidth()
		if s.Changed() {
			d.appAction(inst, map[string]any{"value": *v})
			c.Invalidate()
		}
	case "checkbox":
		v := inst.field(path, "checkbox", n.Checked).Interface().(*bool)
		cb := ui.Checkbox(c, v, n.Label)
		if cb.Changed() {
			d.appAction(inst, map[string]any{"checked": *v})
			c.Invalidate()
		}
	default:
		kit.Muted(c, "?")
	}
}

// renderChildren renders a node's children with stable keys.
func (d *desk) renderChildren(c *ui.Context, inst *appInstance, n *Node, depth int, path string) {
	for i, child := range n.Children {
		c.Key(path + "." + strconv.Itoa(i))
		d.renderNode(c, inst, child, depth+1, path+"."+strconv.Itoa(i))
	}
}

// sendNodeAction maps a node's action to an app action.
func (d *desk) sendNodeAction(inst *appInstance, action any) {
	if m, ok := action.(map[string]any); ok {
		d.appAction(inst, m)
		return
	}
	if action == nil {
		d.appAction(inst, map[string]any{"type": "tap"})
		return
	}
	d.appAction(inst, map[string]any{"value": action})
}

// appAction feeds an action to the app, then refreshes.
func (d *desk) appAction(inst *appInstance, action map[string]any) {
	inst.doAction(action)
}

// knownIcons is the set of valid Lucide icon names.
var knownIcons = map[string]bool{}

func init() {
	for _, n := range icon.Names() {
		knownIcons[n] = true
	}
}

// appColor maps a semantic color name to the theme.
func appColor(t *ui.Theme, name string) ui.Color {
	switch name {
	case "accent":
		return t.Accent
	case "danger":
		return t.Danger
	case "success":
		return t.Success
	case "surface":
		return t.SurfaceHover
	default:
		return t.Surface
	}
}

// textSize maps a size name to DIP and bold.
func textSize(t *ui.Theme, size string) (float32, bool) {
	bold := strings.Contains(size, "bold")
	switch strings.TrimSpace(strings.TrimSuffix(size, "bold")) {
	case "title":
		return t.Rem(1.1), bold
	case "caption":
		return t.Rem(0.72), bold
	default:
		return t.Rem(0.88), bold
	}
}
