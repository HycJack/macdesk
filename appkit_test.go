package main

import (
	"testing"
)

func TestNodeConstructors(t *testing.T) {
	tree := Column(8, Text("7", "title"), Row(4, Button("=", map[string]any{"key": "="}, true)))
	if tree.Type != "column" || tree.Gap != 8 || len(tree.Children) != 2 {
		t.Fatalf("bad root: %+v", tree)
	}
	if tree.Children[0].Text != "7" || tree.Children[0].Size != "title" {
		t.Errorf("bad text: %+v", tree.Children[0])
	}
	btn := tree.Children[1].Children[0]
	if btn.Type != "button" || btn.Label != "=" || !btn.Primary {
		t.Errorf("bad button: %+v", btn)
	}
	if !treeOK(tree) {
		t.Error("good tree rejected")
	}
	bad := Column(0)
	bad.Children = []*Node{bad}
	if treeOK(bad) {
		t.Error("cycle accepted")
	}
	kids := make([]*Node, 501)
	for i := range kids {
		kids[i] = Text("x", "body")
	}
	if treeOK(Column(0, kids...)) {
		t.Error("over-count accepted")
	}
	deep := Text("x", "body")
	for i := 0; i < 13; i++ {
		deep = Column(0, deep)
	}
	if treeOK(deep) {
		t.Error("over-depth accepted")
	}
}

func TestAppSymbols(t *testing.T) {
	syms := appSymbols()
	table := syms["appkit/appkit"]
	for _, name := range []string{"Node", "Column", "Row", "Box", "Text", "Button",
		"Input", "Image", "Progress", "Slider", "Checkbox"} {
		if _, ok := table[name]; !ok {
			t.Errorf("appkit symbol %q missing", name)
		}
	}
	if table["Column"].Type().NumOut() != 1 {
		t.Error("Column must return one node")
	}
}
