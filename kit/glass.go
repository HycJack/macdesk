package kit

import "github.com/egoist/mygo/ui"

// Glass is a translucent surface colour that reads over the wallpaper.
func Glass(t *ui.Theme, alpha float32) ui.Color {
	if t.Dark {
		return ui.Hex("#1e1e20").Alpha(alpha)
	}
	return ui.Hex("#ffffff").Alpha(alpha)
}

// Hairline is the edge colour of glass and of windows.
func Hairline(t *ui.Theme) ui.Color {
	if t.Dark {
		return ui.Hex("#ffffff").Alpha(0.14)
	}
	return ui.Hex("#000000").Alpha(0.12)
}
