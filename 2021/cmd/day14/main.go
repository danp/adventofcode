package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	tmpl := make(map[string]int)
	for _, p := range pairs(lines[0]) {
		tmpl[p]++
	}
	counts := make(map[string]int)
	for _, c := range lines[0] {
		counts[string(c)]++
	}

	rules := make(map[string]string)
	for _, l := range lines[2:] {
		a, b, _ := strings.Cut(l, " -> ")
		rules[a] = b
	}

	fmt.Printf("tmpl: %v\n", tmpl)
	for i := 0; i < 10; i++ {
		tmpl, counts = step(tmpl, counts, rules)
		fmt.Printf("tmpl: %v\n", tmpl)
		fmt.Printf("counts: %v\n", counts)
	}

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

func step(tmpl map[string]int, counts map[string]int, rules map[string]string) (map[string]int, map[string]int) {
	newTmpl := make(map[string]int)
	newCounts := make(map[string]int)
	for tp, tv := range tmpl {
		r := rules[tp]
		newTmpl[string(tp[0])+r] += tv
		newTmpl[r+string(tp[1])] += tv

		newCounts[string(tp[0])] += tv
		newCounts[string(tp[1])] += tv
		newCounts[r] += tv
	}
	return newTmpl, newCounts
}

func pairs(t string) []string {
	var out []string
	for i := 0; i < len(t)-1; i++ {
		out = append(out, t[i:i+2])
	}
	return out
}
