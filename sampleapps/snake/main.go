//go:build ignore

package main

import (
	"math/rand"
	"strconv"

	"appkit"
)

const size = 12

// snakeState is the app's live state.
type snakeState struct {
	snake [][]int
	dir   []int
	food  []int
	score int
	over  bool
}

func Open() any {
	return &snakeState{
		snake: [][]int{{6, 6}, {5, 6}, {4, 6}},
		dir:   []int{1, 0},
		food:  []int{rand.Intn(size), rand.Intn(size)},
	}
}

func View(s any) *appkit.Node {
	st := s.(*snakeState)
	head := appkit.Text("分数 "+strconv.Itoa(st.score), "body")
	if st.over {
		head = appkit.Text("结束！分数 "+strconv.Itoa(st.score), "title bold")
	}
	grid := appkit.Column(1, head, dpad())
	for y := 0; y < size; y++ {
		var row []*appkit.Node
		for x := 0; x < size; x++ {
			row = append(row, appkit.Box(tileColor(st, x, y)))
		}
		grid.Children = append(grid.Children, appkit.Row(1, row...))
	}
	return grid
}

// tileColor picks the color of one grid cell.
func tileColor(st *snakeState, x, y int) string {
	for _, p := range st.snake {
		if p[0] == x && p[1] == y {
			return "accent"
		}
	}
	if st.food[0] == x && st.food[1] == y {
		return "danger"
	}
	return "surface"
}

// dpad is the four direction buttons.
func dpad() *appkit.Node {
	dir := func(dx, dy int) map[string]any {
		return map[string]any{"type": "dir", "dir": []any{float64(dx), float64(dy)}}
	}
	return appkit.Column(2,
		appkit.Row(2, appkit.Button("↑", dir(0, -1), false)),
		appkit.Row(2,
			appkit.Button("←", dir(-1, 0), false),
			appkit.Button("→", dir(1, 0), false),
		),
		appkit.Row(2, appkit.Button("↓", dir(0, 1), false)),
	)
}

func Action(s any, a map[string]any) any {
	st := s.(*snakeState)
	t, _ := a["type"].(string)
	switch t {
	case "dir":
		if d, ok := a["dir"].([]any); ok && len(d) == 2 {
			nd := []int{int(d[0].(float64)), int(d[1].(float64))}
			if nd[0] != -st.dir[0] || nd[1] != -st.dir[1] {
				st.dir = nd
			}
		}
	case "tick":
		step(st)
	case "restart":
		return Open()
	}
	return s
}

// step advances the snake one cell.
func step(st *snakeState) {
	if st.over {
		return
	}
	head := []int{st.snake[0][0] + st.dir[0], st.snake[0][1] + st.dir[1]}
	if head[0] < 0 || head[1] < 0 || head[0] >= size || head[1] >= size {
		st.over = true
		return
	}
	for _, p := range st.snake {
		if p[0] == head[0] && p[1] == head[1] {
			st.over = true
			return
		}
	}
	st.snake = append([][]int{head}, st.snake...)
	if head[0] == st.food[0] && head[1] == st.food[1] {
		st.score++
		st.food = []int{rand.Intn(size), rand.Intn(size)}
	} else {
		st.snake = st.snake[:len(st.snake)-1]
	}
}
