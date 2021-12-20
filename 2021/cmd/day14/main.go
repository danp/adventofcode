package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	tmpl := lines[0]

	rules := make(map[string]string)
	for _, l := range lines[2:] {
		a, b, _ := strings.Cut(l, " -> ")
		rules[a] = b
	}

	fmt.Printf("tmpl: %v\n", tmpl)
	for i := 0; i < 4; i++ {
		tmpl = step(tmpl, rules)
		n := 100
		if n > len(tmpl) {
			n = len(tmpl)
		}
		fmt.Printf("tmpl: %v\n", tmpl[:n])
	}

	counts := make(map[string]int)
	for _, c := range tmpl {
		counts[string(c)]++
	}
	fmt.Printf("counts: %v\n", counts)

	var cs []string
	for c := range counts {
		cs = append(cs, c)
	}
	sort.Slice(cs, func(i, j int) bool {
		return counts[cs[i]] < counts[cs[j]]
	})

	last := cs[len(cs)-1]
	fmt.Printf("most: %v / %v least: %v / %v diff: %v\n", last, counts[last], cs[0], counts[cs[0]], counts[last]-counts[cs[0]])
}

func step(tmpl string, rules map[string]string) string {
	var n string
	for i, p := range pairs(tmpl) {
		if i == 0 {
			n += string(p[0])
		}
		n += rules[p] + string(p[1])
	}
	return n
}

func pairs(t string) []string {
	var out []string
	for i := 0; i < len(t)-1; i++ {
		out = append(out, t[i:i+2])
	}
	return out
}
