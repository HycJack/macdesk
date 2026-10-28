package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/HycJack/macdesk/kit"
	"github.com/egoist/mygo/ui"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// appManifest is the metadata of an installable app.
type appManifest struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	API         int      `json:"api"`
	Entry       string   `json:"entry"`
	Icon        string   `json:"icon"`
	Permissions []string `json:"permissions"`
}

// entry is the manifest's entry file name.
func (m appManifest) entry() string {
	if m.Entry == "" {
		return "main.go"
	}
	return m.Entry
}

// checkManifest validates the required manifest fields.
func checkManifest(m *appManifest) error {
	if !validID(m.ID) {
		return fmt.Errorf("apphost: bad id %q", m.ID)
	}
	if m.Name == "" {
		return errors.New("apphost: name is required")
	}
	if !validVersion(m.Version) {
		return fmt.Errorf("apphost: bad version %q", m.Version)
	}
	if m.API != 1 {
		return fmt.Errorf("apphost: api %d unsupported", m.API)
	}
	return nil
}

// validVersion reports whether v is three safe segments.
func validVersion(v string) bool {
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return false
	}
	for _, p := range parts {
		if p == "" || len(p) > 16 {
			return false
		}
		for _, r := range p {
			if r != '-' && (r < '0' || r > '9') &&
				(r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
				return false
			}
		}
	}
	return true
}

// validID reports whether id is a legal app id.
func validID(id string) bool {
	if id == "" || len(id) > 32 {
		return false
	}
	for _, r := range id {
		if r != '-' && (r < '0' || r > '9') && (r < 'a' || r > 'z') {
			return false
		}
	}
	return true
}

// installedApp is one app under the apps directory.
type installedApp struct {
	manifest appManifest
	dir      string
}

// appHooks are the checked entry points of one app.
type appHooks struct {
	open, view, action, close reflect.Value
}

// hook evaluates one entry point and checks its signature.
func hook(i *interp.Interpreter, name, want string) (reflect.Value, error) {
	v, err := i.Eval(name)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("apphost: %s missing: %w", name, err)
	}
	if got := v.Type().String(); got != want {
		return reflect.Value{}, fmt.Errorf("apphost: %s is %s, want %s", name, got, want)
	}
	return v, nil
}

// loadAppSource compiles app source and checks its hooks.
func loadAppSource(name, src string, perms []string, storageDir string) (*appHooks, error) {
	i := newAppInterpreter(perms, storageDir)
	if _, err := i.Eval(src); err != nil {
		return nil, fmt.Errorf("apphost: %s does not compile: %w", name, err)
	}
	h := &appHooks{}
	var err error
	if h.open, err = hook(i, "main.Open", "func() interface {}"); err != nil {
		return nil, err
	}
	if h.view, err = hook(i, "main.View", "func(interface {}) *main.Node"); err != nil {
		return nil, err
	}
	if h.action, err = hook(i, "main.Action",
		"func(interface {}, map[string]interface {}) interface {}"); err != nil {
		return nil, err
	}
	h.close, _ = hook(i, "main.Close", "func()")
	return h, nil
}

// newAppInterpreter builds a sandboxed interpreter for one app.
func newAppInterpreter(perms []string, storageDir string) *interp.Interpreter {
	i := interp.New(interp.Options{})
	syms := appSymbols()
	if slices.Contains(perms, "storage") {
		for k, v := range storageSymbols(storageDir)["appkit/appkit"] {
			syms["appkit/appkit"][k] = v
		}
	}
	for _, pkg := range []string{"strconv/strconv", "strings/strings", "math/math",
		"math/rand/rand", "sort/sort"} {
		if s, ok := stdlib.Symbols[pkg]; ok {
			syms[pkg] = s
		}
	}
	if err := i.Use(syms); err != nil {
		panic("apphost: " + err.Error())
	}
	return i
}

// appInstance is one open app: its interpreter and live state.
type appInstance struct {
	id       string
	mu       sync.Mutex
	interp   *interp.Interpreter
	hooks    appHooks
	state    reflect.Value
	tree     atomic.Pointer[Node]
	fields   map[string]reflect.Value
	timer    bool
	lastTick time.Time
	lastErr  string
	broken   string
}

// startAppSource opens a fresh instance from app source.
func startAppSource(id, src string, perms []string, storageDir string) (*appInstance, error) {
	hooks, err := loadAppSource(id, src, perms, storageDir)
	if err != nil {
		return nil, err
	}
	a := &appInstance{id: id, hooks: *hooks, fields: map[string]reflect.Value{},
		timer: slices.Contains(perms, "timer"), lastTick: time.Now()}
	st, err := a.call(a.hooks.open)
	if err != nil {
		return nil, err
	}
	a.state = concrete(st)
	a.refresh()
	return a, nil
}

// concrete unwraps an interface-typed reflect value.
func concrete(v reflect.Value) reflect.Value {
	if v.IsValid() && v.Kind() == reflect.Interface {
		return v.Elem()
	}
	return v
}

// call runs one hook under the lock, recovering app panics.
func (a *appInstance) call(fn reflect.Value, args ...reflect.Value) (res reflect.Value, err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.broken != "" {
		return reflect.Value{}, errors.New(a.broken)
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("apphost: %s crashed: %v", a.id, r)
		}
	}()
	start := time.Now()
	out := fn.Call(args)
	if time.Since(start) > 2*time.Second {
		a.broken = "apphost: app call timed out"
	}
	if len(out) > 0 {
		res = out[0]
	}
	return res, nil
}

// refresh recomputes the view tree from the current state.
func (a *appInstance) refresh() {
	res, err := a.call(a.hooks.view, a.state)
	if err != nil {
		a.lastErr = err.Error()
		return
	}
	a.lastErr = ""
	if node, ok := res.Interface().(*Node); ok {
		a.tree.Store(node)
	}
}

// doAction feeds one action to the app, then refreshes.
func (a *appInstance) doAction(action map[string]any) {
	res, err := a.call(a.hooks.action, a.state, reflect.ValueOf(action))
	if err != nil {
		a.lastErr = err.Error()
		return
	}
	a.state = concrete(res)
	a.refresh()
}

// maybeTick fires the timer action on scheduled frames.
func (a *appInstance) maybeTick(c *ui.Context) {
	if !a.timer || a.broken != "" {
		return
	}
	if time.Since(a.lastTick) >= time.Second {
		a.lastTick = a.lastTick.Add(time.Second)
		a.doAction(map[string]any{"type": "tick"})
		c.Invalidate()
	}
	d := time.Until(a.lastTick.Add(time.Second))
	if d <= 0 {
		d = time.Second
	}
	c.After(d)
}

// field returns the persistent binding of one widget node.
func (a *appInstance) field(key, kind string, init any) reflect.Value {
	full := key + "\x00" + kind
	if p, ok := a.fields[full]; ok {
		if p.Elem().Interface() != init {
			p.Elem().Set(reflect.ValueOf(init))
		}
		return p
	}
	p := reflect.New(reflect.TypeOf(init))
	p.Elem().Set(reflect.ValueOf(init))
	a.fields[full] = p
	return p
}

// storageSymbols binds storage helpers to one app directory.
func storageSymbols(dir string) interp.Exports {
	path := filepath.Join(dir, "storage.json")
	return interp.Exports{"appkit/appkit": {
		"StorageGet": reflect.ValueOf(func(key string) string {
			return readStorage(path)[key]
		}),
		"StoragePut": reflect.ValueOf(func(key, value string) {
			m := readStorage(path)
			m[key] = value
			writeStorage(path, m)
		}),
		"StorageDel": reflect.ValueOf(func(key string) {
			m := readStorage(path)
			delete(m, key)
			writeStorage(path, m)
		}),
	}}
}

// readStorage reads one app's storage map.
func readStorage(path string) map[string]string {
	m := map[string]string{}
	b, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(b, &m)
	}
	return m
}

// writeStorage saves one app's storage map.
func writeStorage(path string, m map[string]string) {
	b, _ := json.Marshal(m)
	if err := os.WriteFile(path, b, 0o644); err != nil {
		kit.Log("appstore", "storage", "err", err.Error())
	}
}

// ensureApp starts the instance behind an open app window.
func (d *desk) ensureApp(id string) *appInstance {
	if inst := d.apps[id]; inst != nil {
		return inst
	}
	ia := d.installed[id]
	if ia == nil {
		inst := &appInstance{id: id, lastErr: "apphost: app is gone"}
		d.apps[id] = inst
		return inst
	}
	src, err := os.ReadFile(filepath.Join(ia.dir, ia.manifest.entry()))
	if err != nil {
		inst := &appInstance{id: id, lastErr: err.Error()}
		d.apps[id] = inst
		return inst
	}
	// Storage lives at the app root, shared across versions.
	inst, err := startAppSource(id, string(src), ia.manifest.Permissions,
		filepath.Join(appsRoot, id))
	if err != nil {
		inst = &appInstance{id: id, lastErr: err.Error()}
	}
	d.apps[id] = inst
	return inst
}

// stopApp drops an app instance, running its Close hook.
func (d *desk) stopApp(id string) {
	inst := d.apps[id]
	if inst == nil {
		return
	}
	delete(d.apps, id)
	if inst.hooks.close.IsValid() {
		_, _ = inst.call(inst.hooks.close)
	}
}
