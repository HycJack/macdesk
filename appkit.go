package main

import (
	"reflect"

	"github.com/traefik/yaegi/interp"
)

// Node is one element of an app's view tree.
type Node struct {
	Type     string
	Text     string
	Label    string
	Action   any
	Children []*Node
	Color    string
	Size     string
	Name     string
	Value    float64
	Min      float64
	Max      float64
	Gap      float64
	Checked  bool
	Primary  bool
}

// Column stacks children vertically with a DIP gap.
func Column(gap float64, children ...*Node) *Node {
	return &Node{Type: "column", Gap: gap, Children: children}
}

// Row lays children out horizontally with a DIP gap.
func Row(gap float64, children ...*Node) *Node {
	return &Node{Type: "row", Gap: gap, Children: children}
}

// Box tints a container; childless boxes render as small tiles.
func Box(color string, children ...*Node) *Node {
	return &Node{Type: "box", Color: color, Children: children}
}

// Text shows a string; size is title, body or caption.
func Text(text, size string) *Node { return &Node{Type: "text", Text: text, Size: size} }

// Button sends its action when clicked; primary tints it.
func Button(label string, action any, primary bool) *Node {
	return &Node{Type: "button", Label: label, Action: action, Primary: primary}
}

// Input edits one line; edits arrive as {"text": value}.
func Input(value, placeholder string, action any) *Node {
	return &Node{Type: "input", Text: value, Label: placeholder, Action: action}
}

// Image shows a Lucide icon by name, size DIP.
func Image(name string, size float64) *Node {
	return &Node{Type: "image", Name: name, Value: size}
}

// Progress shows value out of max.
func Progress(value, max float64) *Node {
	return &Node{Type: "progress", Value: value, Max: max}
}

// Slider sends {"value": v} while dragged.
func Slider(value, min, max float64, action any) *Node {
	return &Node{Type: "slider", Value: value, Min: min, Max: max, Action: action}
}

// Checkbox toggles and sends {"checked": b}.
func Checkbox(label string, checked bool, action any) *Node {
	return &Node{Type: "checkbox", Label: label, Checked: checked, Action: action}
}

// appSymbols is the interpreter table exposing the appkit API.
func appSymbols() interp.Exports {
	return interp.Exports{"appkit/appkit": {
		"Node":     reflect.ValueOf((*Node)(nil)),
		"Column":   reflect.ValueOf(Column),
		"Row":      reflect.ValueOf(Row),
		"Box":      reflect.ValueOf(Box),
		"Text":     reflect.ValueOf(Text),
		"Button":   reflect.ValueOf(Button),
		"Input":    reflect.ValueOf(Input),
		"Image":    reflect.ValueOf(Image),
		"Progress": reflect.ValueOf(Progress),
		"Slider":   reflect.ValueOf(Slider),
		"Checkbox": reflect.ValueOf(Checkbox),
	}}
}

// treeOK reports whether the tree fits size and depth rules.
func treeOK(n *Node) bool {
	seen := map[*Node]bool{}
	type frame struct {
		n *Node
		d int
	}
	stack := []frame{{n, 0}}
	count := 0
	for len(stack) > 0 {
		f := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if f.n == nil {
			continue
		}
		if seen[f.n] || f.d > 12 {
			return false
		}
		seen[f.n] = true
		count++
		if count > 500 {
			return false
		}
		for _, c := range f.n.Children {
			stack = append(stack, frame{c, f.d + 1})
		}
	}
	return true
}
