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

// wallpaperPreset is one desktop picture offered by Settings.
type wallpaperPreset struct {
	angle float32
	stops []ui.GradientStop
}

// wallpapers are the desktop pictures, in the order Settings lists them.
var wallpapers = []wallpaperPreset{
	{168, []ui.GradientStop{
		{Color: ui.Hex("#150c26"), At: 0}, {Color: ui.Hex("#3d1a63"), At: 0.32},
		{Color: ui.Hex("#a3246b"), At: 0.62}, {Color: ui.Hex("#e85b45"), At: 0.85},
		{Color: ui.Hex("#f7b267"), At: 1},
	}},
	{160, []ui.GradientStop{
		{Color: ui.Hex("#0d1117"), At: 0}, {Color: ui.Hex("#1d2430"), At: 0.4},
		{Color: ui.Hex("#2f3a4a"), At: 0.72}, {Color: ui.Hex("#4c5b6e"), At: 0.9},
		{Color: ui.Hex("#8b9aad"), At: 1},
	}},
	{170, []ui.GradientStop{
		{Color: ui.Hex("#04203a"), At: 0}, {Color: ui.Hex("#0b4a72"), At: 0.38},
		{Color: ui.Hex("#1b7fae"), At: 0.66}, {Color: ui.Hex("#63c6dc"), At: 0.88},
		{Color: ui.Hex("#cdf1f8"), At: 1},
	}},
	{162, []ui.GradientStop{
		{Color: ui.Hex("#0a1f14"), At: 0}, {Color: ui.Hex("#17452a"), At: 0.36},
		{Color: ui.Hex("#3f8a4d"), At: 0.66}, {Color: ui.Hex("#a8cd5a"), At: 0.88},
		{Color: ui.Hex("#f0f2b0"), At: 1},
	}},
	{174, []ui.GradientStop{
		{Color: ui.Hex("#1b0f2e"), At: 0}, {Color: ui.Hex("#5a1d4d"), At: 0.34},
		{Color: ui.Hex("#c22e5a"), At: 0.64}, {Color: ui.Hex("#f2683c"), At: 0.86},
		{Color: ui.Hex("#ffd08a"), At: 1},
	}},
}

// accentTones are the accent colours Settings offers, light then dark.
var accentTones = []struct{ light, dark ui.Color }{
	{ui.Hex("#0a7cff"), ui.Hex("#0a84ff")},
	{ui.Hex("#8e44ec"), ui.Hex("#bf5af2")},
	{ui.Hex("#e63462"), ui.Hex("#ff2d55")},
	{ui.Hex("#e8442f"), ui.Hex("#ff453a")},
	{ui.Hex("#e8890a"), ui.Hex("#ff9f0a")},
	{ui.Hex("#25a244"), ui.Hex("#30d158")},
	{ui.Hex("#6c6c70"), ui.Hex("#98989d")},
}

// wallpaper paints the desktop picture.
func wallpaper(c *ui.Context, idx int) ui.Element {
	pre := wallpapers[wallIndex(idx)]
	return ui.Box(c).Fill().LinearGradient(ui.LinearGradient{
		Angle: pre.angle, Stops: pre.stops, Oklab: true,
	})
}

// sheet is a layer pinned to every edge of the desktop.
func sheet(c *ui.Context) ui.Element {
	return ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0)
}

// wallIndex keeps a preset number inside the list.
func wallIndex(idx int) int {
	if idx < 0 || idx >= len(wallpapers) {
		return 0
	}
	return idx
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

// accent returns the accent colour of a tone, light or dark.
func accent(tone int, dark bool) ui.Color {
	tone = clampInt(tone, 0, len(accentTones)-1)
	if dark {
		return accentTones[tone].dark
	}
	return accentTones[tone].light
}

// setTheme applies the macos-web palette, light or dark, with a tone.
func setTheme(t *ui.Theme, dark bool, tone int) {
	t.Radius, t.Spacing, t.FontSize = 6, 4, 13
	base := accent(tone, dark)
	if dark {
		t.Dark = true
		t.Background, t.Surface, t.SurfaceHover = ui.Hex("#1e1e1f"), ui.Hex("#262628"), ui.Hex("#2f2f31")
		t.SurfacePressed, t.Border = ui.Hex("#38383a"), ui.Hex("#ffffff").Alpha(0.12)
		t.Text, t.TextMuted = ui.Hex("#f2f2f4"), ui.Hex("#ebebf5").Alpha(0.58)
		t.Danger, t.Warning, t.Success = ui.Hex("#ff453a"), ui.Hex("#ffd60a"), ui.Hex("#32d74b")
		t.Selection, t.Scrollbar = ui.Hex("#787880").Alpha(0.36), ui.Hex("#ffffff").Alpha(0.22)
		t.Accent, t.AccentHover, t.AccentPressed = base, base.Mix(ui.Hex("#ffffff"), 0.16), base.Mix(ui.Hex("#000000"), 0.16)
		t.AccentText = ui.Hex("#ffffff")
		t.Focus = base.Alpha(0.5)
		return
	}
	t.Background, t.Surface, t.SurfaceHover = ui.Hex("#ffffff"), ui.Hex("#f7f7f7"), ui.Hex("#efeff1")
	t.SurfacePressed, t.Border = ui.Hex("#e5e5e8"), ui.Hex("#000000").Alpha(0.085)
	t.Text, t.TextMuted = ui.Hex("#1d1d1f"), ui.Hex("#3c3c43").Alpha(0.62)
	t.Danger, t.Warning, t.Success = ui.Hex("#ff453a"), ui.Hex("#ff9f0a"), ui.Hex("#30d158")
	t.Selection, t.Scrollbar = ui.Hex("#787880").Alpha(0.18), ui.Hex("#000000").Alpha(0.2)
	t.Accent, t.AccentHover, t.AccentPressed = base, base.Mix(ui.Hex("#ffffff"), 0.12), base.Mix(ui.Hex("#000000"), 0.14)
	t.AccentText = ui.Hex("#ffffff")
	t.Focus = base.Alpha(0.45)
}

// clampInt keeps v between lo and hi.
func clampInt(v, lo, hi int) int {
	return max(lo, min(v, hi))
}
