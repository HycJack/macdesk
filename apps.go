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
			Background(t.Surface.Alpha(0.6)).Children(func() {
			desktop.SourceList(c, p.finderGroups, &d.finderIn).FillHeight()
		})
		ui.Column(c).Grow(1).FillHeight().Children(func() {
			tb, _, _ := desktop.UnifiedToolbar(c, p.finderTitle, &d.query,
				desktop.UnifiedToolbarOptions{})
			tb.FillWidth()
			ui.Scroll(c).Fill().Padding(t.Space(3)).Children(func() {
				files := d.finderFiles(p)
				if len(files) == 0 {
					kit.Muted(c, p.finderEmpty).FillWidth().TextAlign(ui.Center)
					return
				}
				ui.Grid(c).Columns(4).Gap(t.Space(3)).Children(func() {
					for _, f := range files {
						c.Key(f.Name)
						d.fileTile(c, f)
					}
				})
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

// fileTile is one document in the Finder grid.
func (d *desk) fileTile(c *ui.Context, f finderFile) {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(f.Name).Column().FillWidth().
		Gap(t.Space(1)).Padding(t.Space(2)).Radius(t.Radius))
	b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
	kit.Clicked(b, sound.Tap)
	b.Children(func() {
		ui.Box(c).Size(t.Space(10), t.Space(10)).Center().Children(func() {
			glyph(c, f.Mark, t.Space(9)).TextColor(t.Accent)
		})
		ui.Text(c, f.Name).FontSize(t.Rem(0.76)).FillWidth().MinWidth(0).
			TextAlign(ui.Center).SingleLine()
		kit.Muted(c, f.Meta).FontSize(t.Rem(0.66)).FillWidth().
			TextAlign(ui.Center).SingleLine()
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
