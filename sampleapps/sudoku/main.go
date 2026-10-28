//go:build ignore

package main

import (
	"strconv"

	"appkit"
)

// puzzles are three base grids, 0 is empty.
var puzzles = [][]int{
	{
		5, 3, 0, 0, 7, 0, 0, 0, 0,
		6, 0, 0, 1, 9, 5, 0, 0, 0,
		0, 9, 8, 0, 0, 0, 0, 6, 0,
		8, 0, 0, 0, 6, 0, 0, 0, 3,
		4, 0, 0, 8, 0, 3, 0, 0, 1,
		7, 0, 0, 0, 2, 0, 0, 0, 6,
		0, 6, 0, 0, 0, 0, 2, 8, 0,
		0, 0, 0, 4, 1, 9, 0, 0, 5,
		0, 0, 0, 0, 8, 0, 0, 7, 9,
	},
	{
		0, 2, 0, 6, 0, 8, 0, 0, 0,
		5, 8, 0, 0, 0, 9, 7, 0, 0,
		0, 0, 0, 0, 4, 0, 0, 0, 0,
		3, 7, 0, 0, 0, 0, 0, 0, 0,
		6, 0, 0, 0, 8, 0, 4, 0, 9,
		0, 0, 9, 0, 6, 0, 0, 5, 1,
		8, 0, 0, 2, 0, 0, 0, 0, 0,
		0, 0, 5, 0, 0, 0, 8, 1, 7,
		0, 0, 4, 0, 9, 7, 0, 6, 0,
	},
	{
		0, 0, 0, 0, 0, 0, 6, 8, 0,
		0, 0, 0, 0, 7, 3, 0, 0, 9,
		3, 0, 9, 0, 0, 0, 0, 4, 5,
		4, 1, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 7, 0, 6, 0, 9, 0, 0,
		0, 0, 0, 8, 0, 0, 0, 1, 2,
		2, 6, 0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 7, 4, 0, 0, 0, 0,
		0, 5, 8, 0, 0, 6, 0, 0, 0,
	},
}

func Open() any {
	return map[string]any{"cells": clone(puzzles[0]), "sel": -1, "puzzle": 0}
}

// clone copies a puzzle grid.
func clone(p []int) []int {
	out := make([]int, len(p))
	copy(out, p)
	return out
}

func View(s any) *appkit.Node {
	m := s.(map[string]any)
	cells := m["cells"].([]int)
	sel := m["sel"].(int)
	grid := appkit.Column(2, appkit.Text("点格子再点数字", "caption"))
	for y := 0; y < 9; y++ {
		var row []*appkit.Node
		for x := 0; x < 9; x++ {
			i := y*9 + x
			label := ""
			if cells[i] > 0 {
				label = strconv.Itoa(cells[i])
			}
			color := "surface"
			if conflict(m, i) {
				color = "danger"
			} else if i == sel {
				color = "accent"
			}
			row = append(row, appkit.Box(color,
				appkit.Button(label, map[string]any{"cell": float64(i)}, false)))
		}
		grid.Children = append(grid.Children, appkit.Row(1, row...))
	}
	var nums []*appkit.Node
	for d := 1; d <= 9; d++ {
		nums = append(nums, appkit.Button(strconv.Itoa(d), map[string]any{"digit": float64(d)}, false))
	}
	grid.Children = append(grid.Children,
		appkit.Row(3, nums...),
		appkit.Button("清除", map[string]any{"digit": 0.0}, false),
	)
	if done(m) {
		grid.Children = append(grid.Children, appkit.Text("完成！", "title bold"))
	}
	return grid
}

func Action(s any, a map[string]any) any {
	m := s.(map[string]any)
	if cell, ok := a["cell"].(float64); ok {
		m["sel"] = int(cell)
	}
	if digit, ok := a["digit"].(float64); ok {
		sel := m["sel"].(int)
		cells := m["cells"].([]int)
		if sel >= 0 && sel < len(cells) {
			cells[sel] = int(digit)
		}
	}
	return m
}

// conflict reports whether a filled cell clashes in row, column or box.
func conflict(m map[string]any, i int) bool {
	cells := m["cells"].([]int)
	v := cells[i]
	if v == 0 {
		return false
	}
	x, y := i%9, i/9
	for k := 0; k < 9; k++ {
		if k != x && cells[y*9+k] == v {
			return true
		}
		if k != y && cells[k*9+x] == v {
			return true
		}
	}
	bx, by := x/3*3, y/3*3
	for dy := 0; dy < 3; dy++ {
		for dx := 0; dx < 3; dx++ {
			j := (by+dy)*9 + bx + dx
			if j != i && cells[j] == v {
				return true
			}
		}
	}
	return false
}

// done reports whether the grid is full and clean.
func done(m map[string]any) bool {
	cells := m["cells"].([]int)
	for i, v := range cells {
		if v == 0 || conflict(m, i) {
			return false
		}
	}
	return true
}
