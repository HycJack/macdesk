//go:build ignore

package main

import (
	"strconv"
	"strings"

	"appkit"
)

func Open() any { return map[string]any{"expr": "", "out": "0"} }

func View(s any) *appkit.Node {
	m := s.(map[string]any)
	rows := [][]string{
		{"7", "8", "9", "÷"}, {"4", "5", "6", "×"},
		{"1", "2", "3", "−"}, {"0", ".", "=", "+"},
	}
	grid := appkit.Column(4, appkit.Text(m["expr"].(string), "caption"),
		appkit.Text(m["out"].(string), "title bold"))
	for _, row := range rows {
		var btns []*appkit.Node
		for _, k := range row {
			btns = append(btns, appkit.Button(k, map[string]any{"key": k}, k == "="))
		}
		grid.Children = append(grid.Children, appkit.Row(4, btns...))
	}
	grid.Children = append(grid.Children, appkit.Button("C", map[string]any{"key": "C"}, false))
	return grid
}

func Action(s any, a map[string]any) any {
	m := s.(map[string]any)
	key, _ := a["key"].(string)
	expr := m["expr"].(string)
	switch key {
	case "C":
		m["expr"], m["out"] = "", "0"
	case "=":
		m["out"] = evaluate(expr)
		m["expr"] = ""
	default:
		m["expr"] = expr + key
	}
	return m
}

// evaluate computes +−×÷ with precedence, "错误" on failure.
func evaluate(expr string) string {
	toks := tokenize(expr)
	if len(toks) == 0 {
		return "0"
	}
	vals, ops := []float64{}, []string{}
	apply := func(op string) bool {
		if len(vals) < 2 {
			return false
		}
		b, a := vals[len(vals)-1], vals[len(vals)-2]
		vals = vals[:len(vals)-2]
		switch op {
		case "+":
			vals = append(vals, a+b)
		case "−":
			vals = append(vals, a-b)
		case "×":
			vals = append(vals, a*b)
		case "÷":
			if b == 0 {
				return false
			}
			vals = append(vals, a/b)
		}
		return true
	}
	for _, tk := range toks {
		if n, err := strconv.ParseFloat(tk, 64); err == nil {
			vals = append(vals, n)
			continue
		}
		for len(ops) > 0 && prec(ops[len(ops)-1]) >= prec(tk) {
			if !apply(ops[len(ops)-1]) {
				return "错误"
			}
			ops = ops[:len(ops)-1]
		}
		ops = append(ops, tk)
	}
	for len(ops) > 0 {
		if !apply(ops[len(ops)-1]) {
			return "错误"
		}
		ops = ops[:len(ops)-1]
	}
	if len(vals) != 1 {
		return "错误"
	}
	return trimNum(vals[0])
}

// tokenize splits an expression into numbers and operators.
func tokenize(expr string) []string {
	var out []string
	num := ""
	for _, r := range expr {
		if strings.ContainsRune("+-×÷", r) {
			if num != "" {
				out = append(out, num)
				num = ""
			}
			out = append(out, string(r))
			continue
		}
		num += string(r)
	}
	if num != "" {
		out = append(out, num)
	}
	return out
}

// prec ranks operator precedence.
func prec(op string) int {
	if op == "×" || op == "÷" {
		return 2
	}
	return 1
}

// trimNum formats a result without trailing zeros.
func trimNum(v float64) string {
	s := strconv.FormatFloat(v, 'f', 6, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" || s == "-" {
		return "0"
	}
	return s
}
