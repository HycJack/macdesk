//go:build ignore

// extract_icons turns the symbol sheet of StarKnightt/macos-web into standalone
// SVGs under appicons/svg. Source: https://github.com/StarKnightt/macos-web (MIT).
//
//	go run tools/extract_icons.go <index.html>
package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// squircle is the app-icon silhouette: |x|^n + |y|^n = 1 at n = 5, splined
// over 64 points, the same outline macos-web builds for its icons.
func squircle(size float64) string {
	const steps = 64
	pts := make([][2]float64, steps+1)
	for i := 0; i <= steps; i++ {
		t := 2 * math.Pi * float64(i) / steps
		c, s := math.Cos(t), math.Sin(t)
		x := math.Copysign(math.Pow(math.Abs(c), 0.4), c)
		y := math.Copysign(math.Pow(math.Abs(s), 0.4), s)
		pts[i] = [2]float64{size / 2 * (1 + x), size / 2 * (1 + y)}
	}
	prec := 2
	if size < 2 {
		prec = 4
	}
	f := func(v float64) string { return strconv.FormatFloat(v, 'f', prec, 64) }
	n := len(pts) - 1
	var b strings.Builder
	b.WriteString("M" + f(pts[0][0]) + " " + f(pts[0][1]))
	for i := 0; i < n; i++ {
		p0, p1 := pts[(i-1+n)%n], pts[i]
		p2, p3 := pts[(i+1)%n], pts[(i+2)%n]
		b.WriteString("C" + f(p1[0]+(p2[0]-p0[0])/6) + " " + f(p1[1]+(p2[1]-p0[1])/6))
		b.WriteString(" " + f(p2[0]-(p3[0]-p1[0])/6) + " " + f(p2[1]-(p3[1]-p1[1])/6))
		b.WriteString(" " + f(p2[0]) + " " + f(p2[1]))
	}
	b.WriteString("Z")
	return b.String()
}

var symbolRE = regexp.MustCompile(`(?s)<symbol id="([^"]+)"[^>]*viewBox="([^"]+)"[^>]*>(.*?)</symbol>`)
var useRE = regexp.MustCompile(`<use href="#([^"]+)"\s*/>`)

// inline swaps <use href="#id"/> for the body of that symbol, which a
// standalone file cannot resolve on its own; nesting stops at depth 4.
func inline(body string, bodies map[string]string, depth int) string {
	if depth > 4 || !strings.Contains(body, "<use ") {
		return body
	}
	return useRE.ReplaceAllStringFunc(body, func(m string) string {
		id := useRE.FindStringSubmatch(m)[1]
		inner, ok := bodies[id]
		if !ok {
			return m
		}
		return inline(inner, bodies, depth+1)
	})
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run tools/extract_icons.go <index.html>")
		os.Exit(2)
	}
	src, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	out := "appicons/svg"
	if err := os.MkdirAll(out, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	clip, unit := squircle(100), squircle(1)
	bodies := map[string]string{}
	for _, m := range symbolRE.FindAllStringSubmatch(string(src), -1) {
		if _, ok := bodies[m[1]]; !ok {
			bodies[m[1]] = m[3]
		}
	}
	seen := map[string]bool{}
	count := 0
	for _, m := range symbolRE.FindAllStringSubmatch(string(src), -1) {
		id, box := m[1], m[2]
		if seen[id] {
			continue
		}
		seen[id] = true
		body := inline(m[3], bodies, 0)
		var b strings.Builder
		b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="` + box + `">`)
		if strings.Contains(body, "url(#sqc)") || strings.Contains(body, "url(#sqr)") {
			b.WriteString(`<defs><clipPath id="sqc"><path d="` + clip + `"/></clipPath>`)
			b.WriteString(`<clipPath id="sqr" clipPathUnits="objectBoundingBox"><path d="` + unit + `"/></clipPath></defs>`)
		}
		b.WriteString(body + `</svg>`)
		if err := os.WriteFile(filepath.Join(out, id+".svg"), []byte(b.String()), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		count++
	}
	fmt.Printf("wrote %d icons\n", count)
}
