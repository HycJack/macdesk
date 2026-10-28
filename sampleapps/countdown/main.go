//go:build ignore

package main

import (
	"strconv"

	"appkit"
)

// countdownState is the app's live state.
type countdownState struct {
	total, left float64
	running     bool
	mins        string
}

func Open() any { return &countdownState{} }

func View(s any) *appkit.Node {
	st := s.(*countdownState)
	label := "开始"
	if st.running {
		label = "暂停"
	}
	return appkit.Column(8,
		appkit.Text(fmtClock(st.left), "title bold"),
		appkit.Progress(st.left, st.total),
		appkit.Input(st.mins, "分钟", map[string]any{"type": "set"}),
		appkit.Row(6,
			appkit.Button(label, map[string]any{"type": "toggle"}, true),
			appkit.Button("重置", map[string]any{"type": "reset"}, false),
		),
	)
}

func Action(s any, a map[string]any) any {
	st := s.(*countdownState)
	t, _ := a["type"].(string)
	switch t {
	case "set":
		text, _ := a["text"].(string)
		st.mins = text
		if mins, err := strconv.ParseFloat(text, 64); err == nil && !st.running {
			st.total, st.left = mins*60, mins*60
		}
	case "toggle":
		if st.total > 0 {
			st.running = !st.running
		}
	case "reset":
		st.running = false
		st.left = st.total
	case "tick":
		if st.running {
			st.left--
			if st.left <= 0 {
				st.running = false
				st.left = 0
			}
		}
	}
	return s
}

// fmtClock renders seconds as mm:ss.
func fmtClock(sec float64) string {
	s := int(sec)
	return pad2(s/60) + ":" + pad2(s%60)
}

// pad2 formats an integer with two digits.
func pad2(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}
