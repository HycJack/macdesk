package main

import (
	"github.com/HycJack/macdesk/kit"
	"github.com/egoist/mygo/ui"
)

// Metrics follow macos-web: a 25 DIP menu bar, 11 DIP windows, a 54 DIP dock.
const (
	menuH     float32 = 25
	dockH     float32 = 86
	dockIcon  float32 = 54
	winRadius float32 = 11
	titleH    float32 = 38
)

// Wallpaper stops, a Sonoma-like sweep from indigo to amber.
var (
	wallA = ui.Hex("#150c26")
	wallB = ui.Hex("#3d1a63")
	wallC = ui.Hex("#a3246b")
	wallD = ui.Hex("#e85b45")
	wallE = ui.Hex("#f7b267")
)

// wallpaper paints the desktop picture.
func wallpaper(c *ui.Context) ui.Element {
	root := ui.Box(c).Fill()
	root.LinearGradient(ui.LinearGradient{
		Angle: 168, Oklab: true,
		Stops: []ui.GradientStop{
			{Color: wallA, At: 0}, {Color: wallB, At: 0.32},
			{Color: wallC, At: 0.62}, {Color: wallD, At: 0.85}, {Color: wallE, At: 1},
		},
	})
	return root
}

// glow is the soft light the wallpaper carries in its lower half.
func glow(c *ui.Context) {
	ui.Box(c).Absolute().Left(0).Right(0).Bottom(0).Height(320).
		LinearGradient(ui.LinearGradient{
			Angle: 180, From: ui.Hex("#ffffff").Alpha(0), To: ui.Hex("#ffffff").Alpha(0.07),
		})
}

// glass is a translucent panel colour over the wallpaper.
func glass(t *ui.Theme, alpha float32) ui.Color { return kit.Glass(t, alpha) }

// hairline is the 0.5 DIP edge of glass and windows.
func hairline(t *ui.Theme) ui.Color { return kit.Hairline(t) }

// setTheme applies the macos-web palette, light or dark.
func setTheme(t *ui.Theme, dark bool) {
	t.Radius, t.Spacing, t.FontSize = 6, 4, 13
	if dark {
		t.Dark = true
		t.Background, t.Surface, t.SurfaceHover = ui.Hex("#1e1e1f"), ui.Hex("#262628"), ui.Hex("#2f2f31")
		t.SurfacePressed, t.Border = ui.Hex("#38383a"), ui.Hex("#ffffff").Alpha(0.12)
		t.Text, t.TextMuted = ui.Hex("#f2f2f4"), ui.Hex("#ebebf5").Alpha(0.58)
		t.Accent, t.AccentHover, t.AccentPressed = ui.Hex("#0a84ff"), ui.Hex("#3395ff"), ui.Hex("#0a6ed1")
		t.AccentText = ui.Hex("#ffffff")
		t.Danger, t.Warning, t.Success = ui.Hex("#ff453a"), ui.Hex("#ffd60a"), ui.Hex("#32d74b")
		t.Selection, t.Scrollbar = ui.Hex("#787880").Alpha(0.36), ui.Hex("#ffffff").Alpha(0.22)
		t.Focus = ui.Hex("#0a84ff").Alpha(0.5)
		return
	}
	t.Background, t.Surface, t.SurfaceHover = ui.Hex("#ffffff"), ui.Hex("#f7f7f7"), ui.Hex("#efeff1")
	t.SurfacePressed, t.Border = ui.Hex("#e5e5e8"), ui.Hex("#000000").Alpha(0.085)
	t.Text, t.TextMuted = ui.Hex("#1d1d1f"), ui.Hex("#3c3c43").Alpha(0.62)
	t.Accent, t.AccentHover, t.AccentPressed = ui.Hex("#0a7cff"), ui.Hex("#1a8bff"), ui.Hex("#0064d8")
	t.AccentText = ui.Hex("#ffffff")
	t.Danger, t.Warning, t.Success = ui.Hex("#ff453a"), ui.Hex("#ff9f0a"), ui.Hex("#30d158")
	t.Selection, t.Scrollbar = ui.Hex("#787880").Alpha(0.18), ui.Hex("#000000").Alpha(0.2)
	t.Focus = ui.Hex("#0a7cff").Alpha(0.45)
}
