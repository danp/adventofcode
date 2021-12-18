package main

import (
	"fmt"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	graph := make(map[string]map[string]struct{})

	for _, l := range lines {
		a, b, _ := strings.Cut(l, "-")
		if _, ok := graph[a]; !ok {
			graph[a] = make(map[string]struct{})
		}
		if _, ok := graph[b]; !ok {
			graph[b] = make(map[string]struct{})
		}
		graph[a][b] = struct{}{}
		graph[b][a] = struct{}{}
	}

	ps := paths(graph, "start", "end", canVisit1)
	fmt.Println(len(ps), "paths between start and end following part 1 rules")

	ps = paths(graph, "start", "end", canVisit2)
	fmt.Println(len(ps), "paths between start and end following part 2 rules")
}

func paths(graph map[string]map[string]struct{}, start, end string, canVisitSmall func([]string, string) bool) [][]string {
	var paths [][]string
	q := [][]string{{start}}
	for len(q) > 0 {
		p := q[0]
		q = q[1:]

		last := p[len(p)-1]
		if last == end {
			paths = append(paths, p)
			continue
		}

		for n := range graph[last] {
			if n == start {
				continue
			}
			if n != end && n >= "a" {
				if !canVisitSmall(p, n) {
					continue
				}
			}
			newp := make([]string, len(p))
			copy(newp, p)
			newp = append(newp, n)
			q = append(q, newp)
		}
	}

	return paths
}

func canVisit1(p []string, n string) bool {
	for _, pp := range p {
		if pp == n {
			return false
		}
	}
	return true
}

func canVisit2(p []string, n string) bool {
	smallVisits := make(map[string]int)
	for _, pp := range p {
		if pp >= "a" && pp != "end" && pp != "start" {
			smallVisits[pp]++
		}
	}
	if smallVisits[n] == 0 {
		return true
	}
	for _, v := range smallVisits {
		if v > 1 {
			return false
		}
	}
	return true
}
