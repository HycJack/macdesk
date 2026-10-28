# macdesk

A macOS-style desktop that runs **Go apps you install at runtime** — an App Store, a dock, windows, a menu bar,
and small apps (calculator, countdown, snake, sudoku) written against a tiny declarative view API.

macdesk is a demo built on [MyGo](https://github.com/egoist/mygo) (native UI) and
[yaegi](https://github.com/traefik/yaegi) (Go interpreter). The desktop is compiled Go; the apps are `.goapp`
packages interpreted inside a sandbox, so installs need no compiler and no restart.

![macdesk desktop](docs/desktop.png)

## App Store

The dock's cart opens the store: installed apps, installable samples, an import field, and install/open/uninstall
buttons. Installing copies the package into `./apps/<id>/<version>/`, opening runs it, uninstalling deletes it.

![App Store](docs/store.png)

An installed app runs in its own window and lives in the dock:

![countdown app](docs/app.png)

## Quick start

```bash
go run .          # or: go build -o macdesk . && ./macdesk
```

Requires Go 1.27+. Installed apps live in `./apps/` (git-ignored). Delete that folder to start clean.

## Writing an app

An app is a folder (or a zip) with `manifest.json` and `main.go`, at most 2 MiB:

```json
{"id": "countdown", "name": "Countdown", "version": "1.0.0", "api": 1, "permissions": ["timer"]}
```

| Field | Rules |
| --- | --- |
| `id` | `[a-z0-9-]{1,32}`, also the storage folder |
| `name` | shown in the store, the window title and the dock |
| `version` | three numeric segments, `1.0.0` |
| `api` | must be `1` |
| `entry` | source file, default `main.go` |
| `permissions` | `timer` for `After`, `storage` for the storage calls |

`main.go` exports four hooks:

```go
package main

import (
	"strconv"

	"appkit"
)

type state struct{ left int }

func Open() any { return &state{left: 60} }

func View(s any) *appkit.Node {
	st := s.(*state)
	return appkit.Column(8,
		appkit.Text("seconds left: "+strconv.Itoa(st.left), "title"),
		appkit.Button("start", map[string]any{"type": "start"}, true),
	)
}

func Action(s any, a map[string]any) any {
	st := s.(*state)
	if a["type"] == "start" {
		st.left = 60
	}
	return s
}

func Close() {} // optional
```

`View` returns a tree of nodes the host renders; `Action` receives `{"type": "tick"}` for timers and touches, slider,
checkbox and input events, and returns the state to keep. See `sampleapps/` for four complete apps.

### appkit

| Call | What it makes |
| --- | --- |
| `Column(gap, children...)` / `Row(gap, children...)` | vertical / horizontal stack, gap in 4 DIP units |
| `Text(text, size)` | `title`, `title bold`, `body`, `caption` |
| `Button(label, action, primary)` | sends `action` when clicked |
| `Input(value, placeholder, action)` | one line, sends `{"text": value}` |
| `Image(lucideName, size)` | Lucide icon |
| `Progress(value, max)` / `Slider(value, min, max, action)` / `Checkbox(label, checked, action)` | values send `{"value"}` / `{"checked"}` |
| `Box(color, children...)` | tinted container, or a small tile when empty |
| `StorageGet` / `StoragePut` / `StorageDel` | per-app key-value JSON, needs `storage` |

Imports are limited to `appkit` plus `strconv`, `strings`, `math`, `math/rand` and `sort`. There are no goroutines,
files or network: with the `timer` permission the host calls `Action` with `{"type": "tick"}` about once a second.

### Pitfalls when writing apps

- Keep state in a struct pointer (`Open` returns `&state{}`). Mixing booleans and maps trips the interpreter.
- Read map values into a local before a `switch` or `if`; do not branch on an indexed expression directly.
- Assign `Children` instead of building huge variadic calls when a row is dynamic.

## How it works

```text
manifest + main.go  ──yaegi──▶  hooks  ──View()──▶  Node tree  ──render──▶  MyGo widgets
        │                          ▲                    │
        └── apps/<id>/<version>/   └──Action(evt)───────┘
```

- `apphost.go` compiles app source in a sandboxed interpreter, validates hooks, and calls them with a timeout.
- `appkit.go` is the node API and its interpreter symbol table; `viewtree.go` renders nodes with MyGo widgets.
- `installer.go` validates manifests and installs, imports (folder or zip) and prunes old versions.
- `appstore.go` is the store UI; `desk.go` is the desktop, windows, dock and menu bar.
- Widgets (`desktop/`, `display/`, `layout/`, `typography/`, `icon/`, `sound/`, `kit/`) come from
  [CrychicUI](https://github.com/ELY-Open-Source/CrychicUI).

## Tests

```bash
go vet ./... && go test ./...
```

The suite drives MyGo's tester: it installs samples, ticks timers, clicks cells, walks the interpreter's state and
checks the store flow end to end. `MACDESK_SHOTS=1 go test -run TestShootDocs .` refreshes `docs/*.png`.

## Credits

MyGo and its Lucide icons, yaegi, and CrychicUI for the macOS-flavoured widgets. MIT licensed, see [LICENSE](LICENSE).
