package main

import (
	"github.com/HycJack/macdesk/desktop"
)

// finderFile is a document shown in the Finder grid.
type finderFile struct{ Name, Mark, Meta string }

// phrases holds every user-facing string of one language.
type phrases struct {
	finder, safari, mail, notes, photos, music, terminal, settings string
	macHD, projects                                                string
	appleAbout, appleSettings, appleHide, appleQuit                string
	fileT, editT, viewT, windowT, helpT                            string
	fileNew, editUndo, editRedo                                    string
	viewIcons, viewList, winMin, winZoom, helpOpen                 string
	closeL, minL, zoomL                                            string
	ccLabel, langToggle, clockFmt                                  string
	ccT                                                            [4]string
	ccVol, ccBri, ccOn, ccOff                                      string
	wifiTitle, wifiEmpty                                           string
	finderTitle, finderEmpty                                       string
	finderGroups                                                   []desktop.SourceGroup
	finderFiles                                                    []finderFile
	noteGroups                                                     []desktop.SourceGroup
	noteBodies                                                     []string
	aboutFmt, aboutCredits, aboutCopyright                         string
	appErrT                                                        string
	storeT, installedT, availableT                                 string
	openL, installL, uninstallL, importL                           string
	importPH, permsT, storeEmpty, storeErrT                        string
}

var zh = phrases{
	finder: "访达", safari: "Safari浏览器", mail: "邮件", notes: "备忘录",
	photos: "照片", music: "音乐", terminal: "终端", settings: "系统设置",
	macHD: "Macintosh HD", projects: "项目文件夹",
	appleAbout: "关于本机", appleSettings: "系统设置", appleHide: "隐藏", appleQuit: "退出 macdesk",
	fileT: "文件", editT: "编辑", viewT: "显示", windowT: "窗口", helpT: "帮助",
	fileNew: "新建访达窗口", editUndo: "撤销", editRedo: "重做",
	viewIcons: "以图标显示", viewList: "以列表显示",
	winMin: "最小化窗口", winZoom: "缩放窗口", helpOpen: "macdesk 帮助",
	closeL: "关闭", minL: "最小化", zoomL: "放大",
	ccLabel: "控制中心", langToggle: "EN", clockFmt: "01-02 15:04",
	ccT:   [4]string{"Wi-Fi", "蓝牙", "专注", "深色模式"},
	ccVol: "声音", ccBri: "显示器", ccOn: "开", ccOff: "关",
	wifiTitle: "网络", wifiEmpty: "未找到网络", finderTitle: "项目", finderEmpty: "没有匹配的项目",
	finderGroups: []desktop.SourceGroup{
		{Title: "个人收藏", Items: []desktop.SourceItem{
			{Name: "AirDrop", Icon: "airdrop"}, {Name: "最近项目", Icon: "clock"},
			{Name: "应用程序", Icon: "grid-2x2", Count: 12}, {Name: "桌面", Icon: "monitor"},
			{Name: "文稿", Icon: "file", Count: 3}, {Name: "下载", Icon: "download", Count: 2},
		}},
		{Title: "位置", Items: []desktop.SourceItem{{Name: "Macintosh HD", Icon: "hard-drive"}}},
	},
	finderFiles: []finderFile{
		{"设计稿.fig", "pen-tool", "昨天 14:20"}, {"预算表.numbers", "table", "周一 09:00"},
		{"路线图.pdf", "file-text", "10-20"}, {"图标.png", "image", "10-18"},
		{"会议记录.md", "file-text", "10-15"}, {"构建产物", "folder", "10-12"},
		{"演示.key", "presentation", "10-10"}, {"数据.csv", "table", "10-08"},
	},
	noteGroups: []desktop.SourceGroup{{Title: "iCloud", Items: []desktop.SourceItem{
		{Name: "会议记录", Icon: "sticky-note"}, {Name: "购物清单", Icon: "shopping-cart"},
		{Name: "灵感", Icon: "lightbulb"},
	}}},
	noteBodies: []string{
		"会议要点：\n- 1500 个组件全部就位\n- 下一步：macOS 桌面演示",
		"购物清单：\n- 咖啡豆\n- USB-C 线\n- 键盘键帽",
		"灵感：\n- 用 ConicGradient 做色轮\n- Dock 放大曲线再弹一点",
	},
	aboutFmt:       "版本 %s",
	aboutCredits:   "界面由 CrychicUI 与 MyGo 在 GPU 上自绘",
	aboutCopyright: "© 2026 macdesk",
	appErrT:        "应用异常",
	storeT:         "App Store", installedT: "已安装", availableT: "可安装",
	openL: "打开", installL: "安装", uninstallL: "卸载", importL: "导入",
	importPH: "输入 .goapp 包路径", permsT: "权限：",
	storeEmpty: "没有可安装的应用", storeErrT: "安装失败",
}

var en = phrases{
	finder: "Finder", safari: "Safari", mail: "Mail", notes: "Notes",
	photos: "Photos", music: "Music", terminal: "Terminal", settings: "System Settings",
	macHD: "Macintosh HD", projects: "Project Folder",
	appleAbout: "About This Mac", appleSettings: "System Settings", appleHide: "Hide", appleQuit: "Quit macdesk",
	fileT: "File", editT: "Edit", viewT: "View", windowT: "Window", helpT: "Help",
	fileNew: "New Finder Window", editUndo: "Undo", editRedo: "Redo",
	viewIcons: "as Icons", viewList: "as List",
	winMin: "Minimize Window", winZoom: "Zoom Window", helpOpen: "macdesk Help",
	closeL: "Close", minL: "Minimize", zoomL: "Zoom",
	ccLabel: "Control Center", langToggle: "中", clockFmt: "Mon Jan 2  15:04",
	ccT:   [4]string{"Wi-Fi", "Bluetooth", "Focus", "Dark Mode"},
	ccVol: "Sound", ccBri: "Display", ccOn: "On", ccOff: "Off",
	wifiTitle: "Networks", wifiEmpty: "No networks found", finderTitle: "Projects", finderEmpty: "No matches",
	finderGroups: []desktop.SourceGroup{
		{Title: "Favorites", Items: []desktop.SourceItem{
			{Name: "AirDrop", Icon: "airdrop"}, {Name: "Recents", Icon: "clock"},
			{Name: "Applications", Icon: "grid-2x2", Count: 12}, {Name: "Desktop", Icon: "monitor"},
			{Name: "Documents", Icon: "file", Count: 3}, {Name: "Downloads", Icon: "download", Count: 2},
		}},
		{Title: "Locations", Items: []desktop.SourceItem{{Name: "Macintosh HD", Icon: "hard-drive"}}},
	},
	finderFiles: []finderFile{
		{"Design.fig", "pen-tool", "Yesterday 14:20"}, {"Budget.numbers", "table", "Mon 09:00"},
		{"Roadmap.pdf", "file-text", "10-20"}, {"Icon.png", "image", "10-18"},
		{"Meeting notes.md", "file-text", "10-15"}, {"Builds", "folder", "10-12"},
		{"Keynote.key", "presentation", "10-10"}, {"Data.csv", "table", "10-08"},
	},
	noteGroups: []desktop.SourceGroup{{Title: "iCloud", Items: []desktop.SourceItem{
		{Name: "Meeting notes", Icon: "sticky-note"}, {Name: "Shopping list", Icon: "shopping-cart"},
		{Name: "Ideas", Icon: "lightbulb"},
	}}},
	noteBodies: []string{
		"Key points:\n- All 1500 components done\n- Next: macOS desktop demo",
		"Shopping list:\n- Coffee beans\n- USB-C cable\n- Keycaps",
		"Ideas:\n- ConicGradient color wheel\n- Bouncier Dock magnification",
	},
	aboutFmt:       "Version %s",
	aboutCredits:   "Drawn on the GPU by CrychicUI and MyGo",
	aboutCopyright: "© 2026 macdesk",
	appErrT:        "App error",
	storeT:         "App Store", installedT: "Installed", availableT: "Available",
	openL: "Open", installL: "Install", uninstallL: "Uninstall", importL: "Import",
	importPH: "Path to a .goapp package", permsT: "Permissions:",
	storeEmpty: "Nothing to install", storeErrT: "Install failed",
}

// tr is the copy of the current language.
func (d *desk) tr() *phrases {
	if d.lang == "en" {
		return &en
	}
	return &zh
}

// toggleLang switches copy and remaps list selections by index.
func (d *desk) toggleLang() {
	from := d.tr()
	if d.lang == "zh" {
		d.lang = "en"
	} else {
		d.lang = "zh"
	}
	to := d.tr()
	d.finderIn = remapSel(from.finderGroups, to.finderGroups, d.finderIn)
	d.noteIn = remapSel(from.noteGroups, to.noteGroups, d.noteIn)
	d.noteText = d.noteBody()
}

// remapSel moves a selection to the same slot across languages.
func remapSel(from, to []desktop.SourceGroup, cur string) string {
	for gi, g := range from {
		for ii, it := range g.Items {
			if it.Name == cur && gi < len(to) && ii < len(to[gi].Items) {
				return to[gi].Items[ii].Name
			}
		}
	}
	return to[0].Items[0].Name
}
