package main

import (
	"fmt"
	"sort"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	type pair []list
	var all []list
	var ps []pair
	var p pair
	for _, l := range lines {
		if l == "" {
			ps = append(ps, p)
			p = nil
			continue
		}
		x, _ := parse(l)
		p = append(p, x)
		all = append(all, x)
	}
	ps = append(ps, p)
	var correct int
	for i, p := range ps {
		c := compare(p[0], p[1])
		if c == -1 {
			correct += i + 1
		}
	}
	fmt.Printf("correct: %v\n", correct)

	d1, _ := parse(`[[2]]`)
	d1.divider = true
	d2, _ := parse(`[[6]]`)
	d2.divider = true
	all = append(all, d1, d2)

	sort.Slice(all, func(i, j int) bool {
		return compare(all[i], all[j]) == -1
	})

	distress := 1
	for i, x := range all {
		if x.divider {
			distress *= i + 1
		}
	}
	fmt.Printf("distress: %v\n", distress)
}

type node any

type list struct {
	nodes   []node
	divider bool
}

type value struct {
	v int
}

func parse(s string) (list, int) {
	var l list
	if len(s) < 1 {
		panic(s)
	}
	if s[0] != '[' {
		panic(s)
	}
	i := 1 // skip opening '['
	for i < len(s) {
		c := s[i]
		switch c {
		case '[':
			sub, n := parse(s[i:])
			l.nodes = append(l.nodes, sub)
			i += n
			continue
		case ']':
			i++
			return l, i
		case ',':
			i++
			continue
		default:
			var ns string
			for _, c := range s[i:] {
				if c < '0' || c > '9' {
					break
				}
				ns += string(c)
			}
			i += len(ns)
			v := value{scaffold.Int(ns)}
			l.nodes = append(l.nodes, v)
		}
	}
	panic("no")
}

func compare(left, right list) int {
	var i int
	for _, ln := range left.nodes {
		if i >= len(right.nodes) {
			return 1
		}
		rn := right.nodes[i]

		lv, lvok := ln.(value)
		rv, rvok := rn.(value)
		if lvok && rvok {
			if lv.v < rv.v {
				return -1
			}
			if lv.v > rv.v {
				return 1
			}
		}

		ll, llok := ln.(list)
		rl, rlok := rn.(list)
		if llok && rlok {
			c := compare(ll, rl)
			if c != 0 {
				return c
			}
		}

		if lvok && rlok {
			if c := compare(list{nodes: []node{lv}}, rl); c != 0 {
				return c
			}
		}
		if rvok && llok {
			if c := compare(ll, list{nodes: []node{rv}}); c != 0 {
				return c
			}
		}
		i++
	}
	if i < len(right.nodes) {
		return -1
	}
	return 0
}
