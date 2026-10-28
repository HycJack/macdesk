package main

import (
	"reflect"
	"strings"
	"testing"
)

// stateField reads an app state field; yaegi prefixes X.
func stateField(inst *appInstance, name string) reflect.Value {
	return inst.state.Elem().FieldByName("X" + name)
}

func loadSample(t *testing.T, id string) *appInstance {
	t.Helper()
	s := sampleByID(id)
	if s == nil {
		t.Fatalf("sample %q missing", id)
	}
	inst, err := startAppSource(id, s.files["main.go"], s.perms, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return inst
}

func TestSamplesLoad(t *testing.T) {
	for _, id := range []string{"calc", "countdown", "snake", "sudoku"} {
		if sampleByID(id) == nil {
			t.Errorf("sample %q not embedded", id)
			continue
		}
		if _, err := startAppSource(id, sampleByID(id).files["main.go"],
			sampleByID(id).perms, t.TempDir()); err != nil {
			t.Errorf("%s does not load: %v", id, err)
		}
	}
}

// findButton locates the first button with a label.
func findButton(n *Node, label string) *Node {
	if n == nil {
		return nil
	}
	if n.Type == "button" && n.Label == label {
		return n
	}
	for _, c := range n.Children {
		if b := findButton(c, label); b != nil {
			return b
		}
	}
	return nil
}

func TestCalcEvaluates(t *testing.T) {
	inst := loadSample(t, "calc")
	for _, key := range []string{"7", "+", "3", "="} {
		inst.doAction(map[string]any{"key": key})
	}
	if !hasText(inst.tree.Load(), "10") {
		t.Error("7+3 did not show 10")
	}
	inst.doAction(map[string]any{"key": "C"})
	if !hasText(inst.tree.Load(), "0") {
		t.Error("clear did not reset")
	}
}

func TestCountdownToggleAndTick(t *testing.T) {
	inst := loadSample(t, "countdown")
	inst.doAction(map[string]any{"type": "set", "text": "1"})
	// Drive the real toggle button action from the tree.
	toggle := findButton(inst.tree.Load(), "开始")
	if toggle == nil {
		t.Fatal("no 开始 button")
	}
	inst.doAction(toggle.Action.(map[string]any))
	if findButton(inst.tree.Load(), "暂停") == nil {
		t.Error("running app did not relabel to 暂停")
	}
	if !stateField(inst, "running").Bool() || stateField(inst, "left").Float() != 60 {
		t.Errorf("running %v left %v", stateField(inst, "running").Bool(),
			stateField(inst, "left").Float())
	}
	inst.doAction(map[string]any{"type": "tick"})
	if stateField(inst, "left").Float() != 59 {
		t.Errorf("tick left %v", stateField(inst, "left").Float())
	}
	// Paused: the label flips back and ticks do nothing.
	inst.doAction(findButton(inst.tree.Load(), "暂停").Action.(map[string]any))
	if stateField(inst, "running").Bool() {
		t.Error("did not pause")
	}
	inst.doAction(map[string]any{"type": "tick"})
	if stateField(inst, "left").Float() != 59 {
		t.Errorf("paused tick moved left to %v", stateField(inst, "left").Float())
	}
}

func TestSnakeMovesAndCollides(t *testing.T) {
	inst := loadSample(t, "snake")
	snakeOf := func() [][]int {
		return stateField(inst, "snake").Interface().([][]int)
	}
	head0 := []int{snakeOf()[0][0], snakeOf()[0][1]}
	// Turn up through the real dpad button, then step.
	up := findButton(inst.tree.Load(), "↑")
	if up == nil {
		t.Fatal("no ↑ button")
	}
	inst.doAction(up.Action.(map[string]any))
	inst.doAction(map[string]any{"type": "tick"})
	if stateField(inst, "over").Bool() {
		t.Fatal("collided without moving into anything")
	}
	head := snakeOf()[0]
	if head[0] != head0[0] || head[1] != head0[1]-1 {
		t.Errorf("head %v, want %v minus one row", head, head0)
	}
	// Force a wall collision one cell from the top edge.
	stateField(inst, "snake").Interface().([][]int)[0] = []int{head[0], 0}
	stateField(inst, "dir").Set(reflect.ValueOf([]int{0, -1}))
	inst.doAction(map[string]any{"type": "tick"})
	if !stateField(inst, "over").Bool() {
		t.Error("wall collision did not end the game")
	}
	if !hasText(inst.tree.Load(), "结束！") {
		t.Error("game over not shown")
	}
}

func TestSudokuConflict(t *testing.T) {
	inst := loadSample(t, "sudoku")
	inst.doAction(map[string]any{"cell": float64(1), "digit": float64(5)})
	if !hasColor(inst.tree.Load(), "danger") {
		t.Error("conflict not marked")
	}
}

// hasText reports whether the tree shows s as a substring.
func hasText(n *Node, s string) bool {
	if n == nil {
		return false
	}
	if strings.Contains(n.Text, s) {
		return true
	}
	for _, c := range n.Children {
		if hasText(c, s) {
			return true
		}
	}
	return false
}

// hasColor reports whether the tree uses a color somewhere.
func hasColor(n *Node, color string) bool {
	if n == nil {
		return false
	}
	if n.Color == color {
		return true
	}
	for _, c := range n.Children {
		if hasColor(c, color) {
			return true
		}
	}
	return false
}
