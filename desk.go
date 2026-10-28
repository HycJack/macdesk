package main

import (
	"maps"
	"slices"
	"strings"

	"github.com/HycJack/macdesk/desktop"
	"github.com/HycJack/macdesk/icon"
	"github.com/HycJack/macdesk/kit"
	"github.com/egoist/mygo/ui"
)

// Screen metrics of the simulated desktop.
const (
	winW   float32 = 1280
	winH   float32 = 820
	winFW  float32 = 620
	winFH  float32 = 400
	aboutW float32 = 360
	aboutH float32 = 280
)

// winState is one simulated window's geometry.
type winState struct {
	open, mini, zoom bool
	x, y             float32
}

// desk is all state of the desktop demo.
type desk struct {
	lang      string
	controls  desktop.Controls
	ccOpen    bool
	apple     bool
	chosen    string
	selHD     bool
	selProj   bool
	iconOpen  int
	wifi      string
	finderIn  string
	query     string
	noteIn    string
	noteText  string
	wins      map[string]*winState
	apps      map[string]*appInstance
	installed map[string]*installedApp
	storeSel  string
	importIn  string
	storeErr  string
	dockOpen  string
	focus     string
	dark      bool
	lpad      bool
	spot      bool
	spotSel   int
	lpadQ     string
	spotQ     string
	notifOpen bool
	notifs    []notif
	trash     []trashItem
	ctxOpen   bool
	ctxX      float32
	ctxY      float32
	ctxApp    string
	dockRight string
	wall      int
	accent    int
	safariIn  string
	mailBox   string
}

func newDesk() *desk {
	d := &desk{lang: "zh", wifi: "Crychic-5G", wins: map[string]*winState{},
		apps: map[string]*appInstance{}}
	d.controls.Dark = true
	d.wins["finder"] = &winState{open: true, x: 96, y: 74}
	d.wins["notes"] = &winState{open: true, x: 700, y: 336}
	d.wins["about"] = &winState{}
	d.wins["store"] = &winState{x: 430, y: 110}
	d.wins["settings"] = &winState{x: 320, y: 130}
	d.wins["trash"] = &winState{x: 430, y: 190}
	d.wins["safari"] = &winState{x: 240, y: 96}
	d.wins["mail"] = &winState{x: 260, y: 120}
	d.wins["photos"] = &winState{x: 300, y: 130}
	d.wins["music"] = &winState{x: 430, y: 150}
	d.wins["terminal"] = &winState{x: 320, y: 170}
	d.finderIn = zh.finderGroups[0].Items[0].Name
	d.noteIn = zh.noteGroups[0].Items[0].Name
	d.noteText = zh.noteBodies[0]
	d.mailBox = zh.mailInbox
	d.scanInstalled()
	return d
}

// view builds the whole desktop from state.
func (d *desk) view(c *ui.Context) {
	d.applyTheme(c)
	p := d.tr()
	d.shortcuts(c)
	root := ui.Box(c).Fill()
	root.Children(func() {
		wallpaper(c, d.wall)
		stage := sheet(c)
		if d.lpad || d.spot {
			stage.Blur(c.Theme().Space(4))
		}
		stage.Children(func() {
			d.icons(c, p)
			d.window(c, p, "finder", p.finderTitle, winFW, winFH, d.finderBody)
			d.window(c, p, "notes", d.noteIn, 560, 380, d.notesBody)
			d.window(c, p, "about", p.appleAbout, aboutW, aboutH, d.aboutBody)
			d.window(c, p, "settings", p.settings, 600, 470, d.settingsBody)
			d.window(c, p, "trash", p.trashT, 440, 380, d.trashBody)
			d.window(c, p, "safari", p.safari, 900, 560, d.safariBody)
			d.window(c, p, "mail", p.mail, 820, 520, d.mailBody)
			d.window(c, p, "photos", p.photos, 760, 520, d.photosBody)
			d.window(c, p, "music", p.music, 420, 480, d.musicBody)
			d.window(c, p, "terminal", p.terminal, 640, 400, d.terminalBody)
			for _, id := range slices.Sorted(maps.Keys(d.installed)) {
				wid := "app:" + id
				if ws := d.wins[wid]; ws == nil || !ws.open {
					continue
				}
				inst := d.ensureApp(id)
				ia := d.installed[id]
				d.window(c, p, wid, ia.manifest.Name, 620, 460, func(cc *ui.Context) {
					d.renderAppView(cc, inst)
				})
			}
			d.window(c, p, "store", p.storeT, 760, 520, d.storeBody)
			d.dockBar(c, p)
		})
		d.menuBar(c, p)
		if d.notifOpen {
			d.notifTray(c, p)
		}
		if d.lpad {
			d.launchpad(c, p)
		}
		if d.spot {
			d.spotlight(c, p)
		}
		if d.ctxOpen {
			if d.ctxApp != "" {
				d.dockMenu(c, p)
			} else {
				d.deskMenu(c, p)
			}
		}
	})
	d.desktopClick(c, root)
}

// shortcuts binds the system keys of the desktop.
func (d *desk) shortcuts(c *ui.Context) {
	c.OnShortcut(ui.Super, ui.KeySpace, func() {
		d.spot = true
		d.spotQ = ""
		d.spotSel = 0
	})
	c.OnShortcut(0, ui.KeyF4, func() {
		d.lpad = !d.lpad
		d.lpadQ = ""
	})
}

// desktopClick opens the desktop menu on a secondary click.
func (d *desk) desktopClick(c *ui.Context, root ui.Element) {
	x, y, over := root.PointerPosition()
	if root.RightClicked() && over {
		d.ctxOpen = true
		d.ctxX, d.ctxY = x, y
		c.Invalidate()
		kit.Log("desktop", "context", "at", int(x), int(y))
	}
}

// SetTheme lasts one frame, so it runs on every build.
func (d *desk) applyTheme(c *ui.Context) {
	t := *c.Theme()
	setTheme(&t, d.controls.Dark, d.accent)
	c.SetTheme(&t)
	d.dark = d.controls.Dark
}

// icons are the desktop shortcuts.
func (d *desk) icons(c *ui.Context, p *phrases) {
	t := c.Theme()
	col := ui.Column(c).Absolute().Top(menuH + t.Space(4)).Left(t.Space(4)).Gap(t.Space(1.5))
	col.Children(func() {
		before := d.iconOpen
		desktop.DesktopIcon(c, p.macHD, "f-hd", &d.selHD, &d.iconOpen)
		desktop.DesktopIcon(c, p.projects, "f-folder-blue", &d.selProj, &d.iconOpen)
		if d.iconOpen > before {
			d.openID("finder")
		}
	})
}

// window floats one app window with lights and dragging.
func (d *desk) window(c *ui.Context, p *phrases, id, title string, w, h float32, body func(*ui.Context)) {
	ws := d.wins[id]
	if !ws.open || ws.mini {
		return
	}
	t := c.Theme()
	surf := ui.Box(c)
	if ws.zoom {
		surf = surf.Absolute().Top(menuH).Left(0).Right(0).Bottom(dockH)
	} else {
		surf = surf.Absolute().Top(ws.y).Left(ws.x).Size(w, h)
	}
	shadow := kit.Shade(0.32)
	if d.focus == id {
		shadow = kit.Shade(0.52)
	}
	surf.Clip().Radius(winRadius).Background(t.Background).Border(1, hairline(t)).
		Shadow(0, 22, 56, -8, shadow).Role(ui.RoleGroup).Label(title)
	surf.Children(func() {
		d.titleBar(c, p, id, title)
		ui.Box(c).Fill().Children(func() { body(c) })
	})
}

// titleBar holds the lights, the title and the drag surface.
func (d *desk) titleBar(c *ui.Context, p *phrases, id, title string) {
	ws := d.wins[id]
	t := c.Theme()
	bar := ui.Row(c).FillWidth().Height(titleH).Gap(t.Space(3)).AlignItems(ui.Center).
		Padding(0, t.Space(3)).Background(t.Surface).
		BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
	bar.Children(func() {
		_, close, mini, zoom := desktop.TrafficLights(c, desktop.TrafficLightsOptions{
			Close: p.closeL, Minimize: p.minL, Zoom: p.zoomL,
		})
		if close.Clicked() {
			ws.open = false
			if strings.HasPrefix(id, "app:") {
				d.stopApp(strings.TrimPrefix(id, "app:"))
			}
			bar.SetChanged()
			kit.Log("macdesk", "close", "win", id)
		}
		if mini.Clicked() {
			ws.mini = true
			bar.SetChanged()
			kit.Log("macdesk", "minimize", "win", id)
		}
		if zoom.Clicked() {
			ws.zoom = !ws.zoom
			bar.SetChanged()
			kit.Log("macdesk", "zoom", "win", id)
		}
		ui.Text(c, title).FontSize(t.Rem(0.95)).Bold().Grow(1).MinWidth(0).
			TextAlign(ui.Center).SingleLine().Unselectable()
		ui.Row(c).Width(t.Space(13)).Shrink(0)
	})
	if bar.Pressed() {
		d.focus = id
	}
	d.drag(c, ws, bar)
}

// drag moves the window while its title bar is pressed.
func (d *desk) drag(c *ui.Context, ws *winState, bar ui.Element) {
	if ws.zoom {
		return
	}
	dx, dy, ok := bar.Dragged()
	if !ok {
		return
	}
	if dx == 0 && dy == 0 {
		return
	}
	ws.x = clamp(ws.x+dx, -winFW+160, winW-160)
	ws.y = clamp(ws.y+dy, menuH, winH-dockH-40)
	c.Invalidate()
}

// dockBar centers the Dock on the bottom edge.
func (d *desk) dockBar(c *ui.Context, p *phrases) {
	t := c.Theme()
	strip := ui.Row(c).Absolute().Left(0).Right(0).Bottom(t.Space(1.5)).Height(dockH).
		Justify(ui.Center).AlignItems(ui.End)
	strip.Children(func() {
		apps := []desktop.DockApp{
			{Name: p.finder, Art: "app-finder", Running: d.wins["finder"].open},
			{Name: p.launchpadT, Art: "app-launchpad", Running: d.lpad},
			{Name: p.safari, Art: "app-safari", Running: d.wins["safari"].open},
			{Name: p.mail, Art: "app-mail", Badge: 3, Running: d.wins["mail"].open},
			{Name: p.notes, Art: "app-notes", Running: d.wins["notes"].open},
			{Name: p.photos, Art: "app-photos", Running: d.wins["photos"].open},
			{Name: p.music, Art: "app-music", Running: d.wins["music"].open},
			{Name: p.terminal, Art: "app-terminal", Running: d.wins["terminal"].open},
			{Name: p.settings, Art: "app-settings", Running: d.wins["settings"].open},
			{Name: p.storeT, Art: "app-appstore", Running: d.wins["store"].open},
		}
		for _, id := range sortedApps(d.installed) {
			apps = append(apps, desktop.DockApp{Name: d.installed[id].manifest.Name,
				Art:     appArt(id, d.installed[id].manifest.Icon),
				Running: d.wins["app:"+id] != nil && d.wins["app:"+id].open})
		}
		apps = append(apps, desktop.DockApp{Name: p.trashT, Art: trashArt(d.trash),
			Running: d.wins["trash"].open, Separator: true})
		dock := desktop.Dock(c, apps, &d.dockOpen, &d.dockRight)
		if dock.Changed() {
			if d.dockRight != "" {
				d.openDockMenu(c, p, dock)
			} else if d.dockOpen != "" {
				d.openApp(p, d.dockOpen)
				d.dockOpen = ""
			}
		}
	})
}

// openApp opens the window of a dock app by name.
func (d *desk) openApp(p *phrases, name string) {
	id := ""
	switch name {
	case p.finder:
		id = "finder"
	case p.notes:
		id = "notes"
	case p.storeT:
		id = "store"
	case p.settings:
		id = "settings"
	case p.trashT:
		id = "trash"
	case p.safari:
		id = "safari"
	case p.mail:
		id = "mail"
	case p.photos:
		id = "photos"
	case p.music:
		id = "music"
	case p.terminal:
		id = "terminal"
	case p.launchpadT:
		d.lpad = !d.lpad
		return
	default:
		for aid, ia := range d.installed {
			if ia.manifest.Name == name {
				id = "app:" + aid
				break
			}
		}
	}
	if id == "" {
		kit.Log("macdesk", "dock", "app", name)
		return
	}
	d.openID(id)
}

// openID opens or restores a window and focuses it.
func (d *desk) openID(id string) {
	ws := d.wins[id]
	if ws == nil {
		ws = &winState{x: d.nextWinX(), y: d.nextWinY()}
		d.wins[id] = ws
	}
	ws.open, ws.mini = true, false
	d.focus = id
	kit.Log("macdesk", "open", "win", id)
}

// nextWinX cascades new windows away from the corner.
func (d *desk) nextWinX() float32 {
	return 150 + float32(len(d.wins)%6)*52
}

// nextWinY cascades new windows away from the corner.
func (d *desk) nextWinY() float32 {
	return 70 + float32(len(d.wins)%6)*40
}

// glyph is a Lucide icon size DIP square.
func glyph(c *ui.Context, name string, size float32) ui.Element {
	return ui.Icon(c, icon.Must(name)).Size(size, size).Shrink(0)
}

// unit is 1 for true and 0 for false.
func unit(v bool) float32 {
	if v {
		return 1
	}
	return 0
}

// clamp keeps v in lo..hi.
func clamp(v, lo, hi float32) float32 { return min(max(v, lo), hi) }
