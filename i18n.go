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
	searchPH, backL, fwdL                                          string
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
	storeTag, storeEntryT, storeAPIT, storeLimitT                  string
	moreT                                                          string
	launchpadT, trashT, spotT, spotPH, notifT                      string
	spotApps, spotActions, spotFiles, spotEmpty                    string
	actDark, actWall, actLpad, actNotif, actTrash, actAbout        string
	wallNames, accentNames                                         []string
	setAppearance, setLight, setDark, setWallpaper, setAccent      string
	setAbout, setVersion                                           string
	trashEmptyT, trashEmptyHint, trashClear, trashHint             string
	notifClear, notifNone, notifInstall, notifRemove               string
	calWeek                                                        [7]string
	today, weatherCity, weatherState                               string
	dockOpenL, dockCloseL, dockQuitL                               string
	ctxNewFolder, ctxWall, ctxSpot, ctxLpad, ctxTrash, ctxAbout    string
	favT, safariStart, safariPH                                    string
	mailBoxes, mailInbox, mailSent, mailDrafts                     string
	musicNow, playL                                                string
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
	searchPH: "搜索", backL: "后退", fwdL: "前进",
	finderGroups: []desktop.SourceGroup{
		{Title: "个人收藏", Items: []desktop.SourceItem{
			{Name: "AirDrop", Icon: "g-airdrop"}, {Name: "最近项目", Icon: "g-clock-a"},
			{Name: "应用程序", Icon: "g-apps", Count: 12}, {Name: "桌面", Icon: "g-desktop"},
			{Name: "文稿", Icon: "g-docs", Count: 3}, {Name: "下载", Icon: "g-download", Count: 2},
		}},
		{Title: "位置", Items: []desktop.SourceItem{{Name: "Macintosh HD", Icon: "g-hd"}}},
	},
	finderFiles: []finderFile{
		{"设计稿.fig", "f-img", "昨天 14:20"}, {"预算表.numbers", "f-code", "周一 09:00"},
		{"路线图.pdf", "f-txt", "10-20"}, {"图标.png", "f-img", "10-18"},
		{"会议记录.md", "f-txt", "10-15"}, {"构建产物", "f-folder", "10-12"},
		{"演示.key", "f-pdf", "10-10"}, {"数据.csv", "f-code", "10-08"},
	},
	noteGroups: []desktop.SourceGroup{{Title: "iCloud", Items: []desktop.SourceItem{
		{Name: "会议记录", Icon: "g-doc-g"}, {Name: "购物清单", Icon: "g-checklist"},
		{Name: "灵感", Icon: "g-star"},
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
	storeTag:    "用 Go 写的小软件，装完即用：界面由系统绘制，脚本在沙箱里运行。",
	storeEntryT: "入口", storeAPIT: "接口", storeLimitT: "上限",
	moreT:      "这里还有",
	launchpadT: "启动台", trashT: "废纸篓", spotT: "聚焦搜索", spotPH: "搜索应用、文件与操作",
	notifT:   "通知中心",
	spotApps: "应用程序", spotActions: "操作", spotFiles: "文件", spotEmpty: "没有匹配的结果",
	actDark: "切换深色模式", actWall: "更换壁纸", actLpad: "打开启动台",
	actNotif: "通知中心", actTrash: "清空废纸篓", actAbout: "关于本机",
	wallNames:     []string{"紫霞", "石墨", "海洋", "森林", "日落"},
	accentNames:   []string{"蓝色", "紫色", "粉色", "红色", "橙色", "绿色", "石墨"},
	setAppearance: "外观", setLight: "浅色", setDark: "深色", setWallpaper: "壁纸",
	setAccent: "强调色", setAbout: "关于", setVersion: "版本",
	trashEmptyT: "废纸篓是空的", trashEmptyHint: "卸载的应用会先放到这里",
	trashClear: "清空", trashHint: "已卸载的应用",
	notifClear: "全部清除", notifNone: "没有新通知",
	notifInstall: "已安装", notifRemove: "已卸载",
	calWeek: [7]string{"一", "二", "三", "四", "五", "六", "日"},
	today:   "今天", weatherCity: "上海", weatherState: "多云 24°",
	dockOpenL: "打开", dockCloseL: "关闭窗口", dockQuitL: "退出",
	ctxNewFolder: "新建文件夹", ctxWall: "更换壁纸", ctxSpot: "聚焦搜索",
	ctxLpad: "启动台", ctxTrash: "清空废纸篓", ctxAbout: "关于本机",
	favT: "个人收藏", safariStart: "起始页", safariPH: "搜索或输入网站名称",
	mailBoxes: "邮箱", mailInbox: "收件箱", mailSent: "已发送", mailDrafts: "草稿",
	musicNow: "正在播放", playL: "播放",
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
	searchPH: "Search", backL: "Back", fwdL: "Forward",
	finderGroups: []desktop.SourceGroup{
		{Title: "Favorites", Items: []desktop.SourceItem{
			{Name: "AirDrop", Icon: "g-airdrop"}, {Name: "Recents", Icon: "g-clock-a"},
			{Name: "Applications", Icon: "g-apps", Count: 12}, {Name: "Desktop", Icon: "g-desktop"},
			{Name: "Documents", Icon: "g-docs", Count: 3}, {Name: "Downloads", Icon: "g-download", Count: 2},
		}},
		{Title: "Locations", Items: []desktop.SourceItem{{Name: "Macintosh HD", Icon: "g-hd"}}},
	},
	finderFiles: []finderFile{
		{"Design.fig", "f-img", "Yesterday 14:20"}, {"Budget.numbers", "f-code", "Mon 09:00"},
		{"Roadmap.pdf", "f-txt", "10-20"}, {"Icon.png", "f-img", "10-18"},
		{"Meeting notes.md", "f-txt", "10-15"}, {"Builds", "f-folder", "10-12"},
		{"Keynote.key", "f-pdf", "10-10"}, {"Data.csv", "f-code", "10-08"},
	},
	noteGroups: []desktop.SourceGroup{{Title: "iCloud", Items: []desktop.SourceItem{
		{Name: "Meeting notes", Icon: "g-doc-g"}, {Name: "Shopping list", Icon: "g-checklist"},
		{Name: "Ideas", Icon: "g-star"},
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
	storeTag:    "Small programs written in Go: the system draws the views, the script stays sandboxed.",
	storeEntryT: "Entry", storeAPIT: "API", storeLimitT: "Limit",
	moreT:      "More in the store",
	launchpadT: "Launchpad", trashT: "Trash", spotT: "Spotlight", spotPH: "Search apps, files and actions",
	notifT:   "Notification Centre",
	spotApps: "Applications", spotActions: "Actions", spotFiles: "Documents", spotEmpty: "No matches",
	actDark: "Toggle Dark Mode", actWall: "Change Wallpaper", actLpad: "Open Launchpad",
	actNotif: "Notification Centre", actTrash: "Empty Trash", actAbout: "About This Mac",
	wallNames:     []string{"Aurora", "Graphite", "Ocean", "Forest", "Sunset"},
	accentNames:   []string{"Blue", "Purple", "Pink", "Red", "Orange", "Green", "Graphite"},
	setAppearance: "Appearance", setLight: "Light", setDark: "Dark", setWallpaper: "Wallpaper",
	setAccent: "Accent colour", setAbout: "About", setVersion: "Version",
	trashEmptyT: "Trash is empty", trashEmptyHint: "Uninstalled apps land here first",
	trashClear: "Empty", trashHint: "Uninstalled apps",
	notifClear: "Clear All", notifNone: "No new notifications",
	notifInstall: "Installed", notifRemove: "Uninstalled",
	calWeek: [7]string{"M", "T", "W", "T", "F", "S", "S"},
	today:   "Today", weatherCity: "Shanghai", weatherState: "Cloudy 24°",
	dockOpenL: "Open", dockCloseL: "Close Window", dockQuitL: "Quit",
	ctxNewFolder: "New Folder", ctxWall: "Change Wallpaper", ctxSpot: "Spotlight",
	ctxLpad: "Launchpad", ctxTrash: "Empty Trash", ctxAbout: "About This Mac",
	favT: "Favourites", safariStart: "Start Page", safariPH: "Search or enter website name",
	mailBoxes: "Mailboxes", mailInbox: "Inbox", mailSent: "Sent", mailDrafts: "Drafts",
	musicNow: "Now Playing", playL: "Play",
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
