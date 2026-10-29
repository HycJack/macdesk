package main

import (
	"fmt"
	"strings"

	"github.com/HycJack/macdesk/desktop"
	"github.com/HycJack/macdesk/kit"
	"github.com/HycJack/macdesk/sound"
	"github.com/egoist/mygo/ui"
)

// finderBody is the Finder window: sidebar, toolbar and file grid.
func (d *desk) finderBody(c *ui.Context) {
	p := d.tr()
	t := c.Theme()
	ui.Row(c).Fill().Children(func() {
		ui.Box(c).Width(t.Space(44)).FillHeight().
			Background(t.Surface.Alpha(0.6)).
			BorderWidth(0, 0, 0, 1).BorderColor(hairline(t)).Children(func() {
			desktop.SourceList(c, p.finderGroups, &d.finderIn).FillHeight()
		})
		ui.Column(c).Grow(1).FillHeight().Children(func() {
			tb, _, _ := desktop.UnifiedToolbar(c, p.finderTitle, &d.query,
				desktop.UnifiedToolbarOptions{Placeholder: p.searchPH, Back: p.backL, Forward: p.fwdL, Compact: true})
			tb.FillWidth()
			files := d.finderFiles(p)
			ui.Scroll(c).Fill().Padding(t.Space(3)).Children(func() {
				if len(files) == 0 {
					kit.Muted(c, p.finderEmpty).FillWidth().TextAlign(ui.Center)
					return
				}
				d.finderGrid(c, files)
			})
			ui.Row(c).FillWidth().Height(t.Space(6)).
				Padding(0, t.Space(3)).AlignItems(ui.Center).
				Background(t.Surface.Alpha(0.7)).
				BorderWidth(1, 0, 0, 0).BorderColor(hairline(t)).Children(func() {
				kit.Muted(c, fmt.Sprintf(p.itemsFmt, len(files))).FontSize(t.Rem(0.66))
			})
		})
	})
}

// finderFiles filters the sample files by the search query.
func (d *desk) finderFiles(p *phrases) []finderFile {
	if d.query == "" {
		return p.finderFiles
	}
	var out []finderFile
	for _, f := range p.finderFiles {
		if strings.Contains(f.Name, d.query) {
			out = append(out, f)
		}
	}
	return out
}

// finderCols is how many tiles fit across the icon view.
func (d *desk) finderCols(t *ui.Theme) int {
	w := float32(winFW)
	if d.wins["finder"].zoom {
		w = winW
	}
	avail := w - t.Space(44) - t.Space(8)
	return min(max(int(avail/t.Space(30)), 2), 8)
}

// finderGrid lays the tiles out in rows of whatever fits.
func (d *desk) finderGrid(c *ui.Context, files []finderFile) {
	t := c.Theme()
	cols := d.finderCols(t)
	ui.Column(c).FillWidth().Gap(t.Space(1.5)).Children(func() {
		for start := 0; start < len(files); start += cols {
			end := min(start+cols, len(files))
			ui.Row(c).FillWidth().Gap(t.Space(1.5)).Children(func() {
				for _, f := range files[start:end] {
					c.Key(f.Name)
					d.fileTile(c, f)
				}
				for i := end - start; i < cols; i++ {
					ui.Box(c).Grow(1)
				}
			})
		}
	})
}

// fileTile is one document in the Finder grid.
func (d *desk) fileTile(c *ui.Context, f finderFile) {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(f.Name).Column().Width(t.Space(30)).Shrink(0).
		MinWidth(0).Clip().Gap(t.Space(1)).Padding(t.Space(2)).Radius(t.Radius))
	b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
	kit.Clicked(b, sound.Tap)
	b.Children(func() {
		ui.Box(c).Width(t.Space(13)).Height(t.Space(12)).Center().Children(func() {
			desktop.Art(c, f.Mark, t.Space(11))
		})
		ui.Text(c, f.Name).FontSize(t.Rem(0.78)).FillWidth().MinWidth(0).
			TextAlign(ui.Center).SingleLine().MaxLines(1).Ellipsis("…")
		kit.Muted(c, f.Meta).FontSize(t.Rem(0.66)).FillWidth().MinWidth(0).
			TextAlign(ui.Center).SingleLine().MaxLines(1).Ellipsis("…")
	})
}

// notesBody is the Notes window: note list and the text.
func (d *desk) notesBody(c *ui.Context) {
	p := d.tr()
	t := c.Theme()
	ui.Row(c).Fill().Children(func() {
		ui.Box(c).Width(t.Space(40)).FillHeight().
			Background(t.Surface.Alpha(0.6)).Children(func() {
			sl := desktop.SourceList(c, p.noteGroups, &d.noteIn).FillHeight()
			if sl.Changed() {
				d.noteText = d.noteBody()
			}
		})
		ui.Column(c).Grow(1).FillHeight().Padding(t.Space(3)).Gap(t.Space(2)).Children(func() {
			ui.Text(c, d.noteIn).FontSize(t.Rem(1.05)).Bold().FillWidth().SingleLine()
			ui.TextArea(c, &d.noteText).Fill().Grow(1)
		})
	})
}

// noteBody is the text of the open note.
func (d *desk) noteBody() string {
	p := d.tr()
	for j, it := range p.noteGroups[0].Items {
		if it.Name == d.noteIn {
			return p.noteBodies[j]
		}
	}
	return p.noteBodies[0]
}

// aboutBody is the About window content.
func (d *desk) aboutBody(c *ui.Context) {
	p := d.tr()
	t := c.Theme()
	ui.Column(c).Fill().Center().Gap(t.Space(2)).Padding(t.Space(6)).Children(func() {
		ui.Box(c).Size(t.Space(16), t.Space(16)).Center().Radius(t.Radius * 3).
			Background(t.Accent.Alpha(0.15)).TextColor(t.Accent).
			Children(func() { glyph(c, "sparkles", t.Space(9)) })
		ui.Text(c, "macdesk").FontSize(t.Rem(1.2)).Bold()
		ui.Text(c, fmt.Sprintf(p.aboutFmt, "v0.5.5")).Font("monospace").
			FontSize(t.Rem(0.76)).TextColor(t.TextMuted)
		kit.Muted(c, p.aboutCredits).TextAlign(ui.Center).Wrap().MaxWidth(t.Space(70))
		ui.Text(c, p.aboutCopyright).FontSize(t.Rem(0.68)).TextColor(t.TextMuted)
	})
}
