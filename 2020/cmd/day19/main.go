package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	graph := make(map[string]rule)
	for _, l := range lines {
		if l == "" {
			break
		}
		if l[0] == '#' {
			continue
		}
		colon := strings.Index(l, ":")
		if colon == -1 {
			panic("no")
		}

		n := l[:colon]
		l = l[colon+2:]

		nd := rule{
			name: n,
		}
		if l[0] == '"' && l[len(l)-1] == '"' {
			nd.lit = l[1 : len(l)-1]
		} else {
			var seq []string
			for _, f := range strings.Fields(l) {
				if f == "|" {
					nd.seqs = append(nd.seqs, seq)
					seq = nil
					continue
				}
				seq = append(seq, f)
			}
			nd.seqs = append(nd.seqs, seq)
		}

		graph[n] = nd
	}

	rs := rules{g: graph}

	var st int
	var valid int
	for _, l := range lines {
		if l == "" {
			if st == 1 {
				break
			}
			st++
			continue
		}

		if st == 0 {
			continue
		}

		m := rs.match(l)
		if m {
			valid++
		}
		fmt.Println(l, m)
	}

	fmt.Println(valid, "are valid")

	rs.g["8"] = rule{name: "8", seqs: [][]string{{"42"}, {"42", "8"}}}
	rs.g["11"] = rule{name: "11", seqs: [][]string{{"42", "31"}, {"42", "11", "31"}}}

	st = 0
	valid = 0
	for _, l := range lines {
		if l == "" {
			if st == 1 {
				break
			}
			st++
			continue
		}

		if st == 0 {
			continue
		}

		m := rs.match(l)
		if m {
			valid++
		}
		fmt.Println(l, m)
	}

	fmt.Println(valid, "are valid")
}

type rules struct {
	g map[string]rule
}

func (r rules) match(s string) bool {
	matches := r.g["0"].match(r.g, s)
	fmt.Println(s, matches)
	if len(matches) == 0 {
		return false
	}
	for _, m := range matches {
		if len(m) == len(s) {
			return true
		}
	}

	return false
}

type rule struct {
	name string
	// if len > 1 it's an OR
	seqs [][]string
	lit  string
}

func (r rule) match(graph map[string]rule, s string) []string {
	if len(s) == 0 {
		return nil
	}

	if r.lit == string(s[0]) {
		return []string{r.lit}
	}

	if len(r.seqs) > 0 && r.seqs[len(r.seqs)-1][len(r.seqs[len(r.seqs)-1])-1] == r.name {
		var out []string
		var start int
		rule := graph[r.seqs[0][0]]
		for {
			matches := rule.match(graph, s[start:])
			if len(matches) > 0 {
				for _, m := range matches {
					out = append(out, s[:start]+m)
				}
				start += len(matches[len(matches)-1])
			} else {
				break
			}
		}

		return out
	}

	for _, rs := range r.seqs {
		cons := map[int]bool{0: true}
		allmatch := true
		for _, sn := range rs {
			newcons := make(map[int]bool)
			for st := range cons {
				pm := graph[sn].match(graph, s[st:])
				if len(pm) == 0 {
					continue
				}
				for _, pmm := range pm {
					newcons[st+len(pmm)] = true
				}
			}
			cons = newcons
			if len(newcons) == 0 {
				allmatch = false
				break
			}
		}
		if allmatch {
			out := make([]string, 0, len(cons))
			for st := range cons {
				out = append(out, s[:st])
			}
			sort.Slice(out, func(i, j int) bool { return len(out[i]) < len(out[j]) })
			return out
		}
	}

	return nil
}
