package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/danp/adventofcode/scaffold"
	"golang.org/x/exp/maps"
)

func main() {
	lines := scaffold.Lines()

	monkeys := make(map[int]monkey)
	var mid int
	var c monkey
	divprod := 1
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "Monkey") {
			if len(c.items) > 0 {
				monkeys[mid] = c
			}
			mid = scaffold.Int(l[len("Monkey ") : len(l)-1])
			c = monkey{}
		}
		if strings.HasPrefix(l, "Starting items") {
			is := strings.Split(l[len("Starting items: "):], ", ")
			for _, i := range scaffold.Ints(is) {
				c.items = append(c.items, item{i, nil, nil})
			}
		}
		if strings.HasPrefix(l, "Operation: new = old ") {
			ops := l[len("Operation: new = old "):]
			if ops[0] == '*' {
				if ops[2:] == "old" {
					c.op = func(it item) item {
						res := it.n * it.n
						it.ops = append(it.ops, fmt.Sprintf("m%v: old(%v) * old(%v) = %v", mid, it.n, it.n, res))
						it.prevs = append(it.prevs, it.n)
						it.n = res
						return it
					}
				} else {
					n := scaffold.Int(ops[2:])
					c.op = func(it item) item {
						res := it.n * n
						it.ops = append(it.ops, fmt.Sprintf("m%v: old(%v) * %v = %v", mid, it.n, n, res))
						it.prevs = append(it.prevs, it.n)
						it.n = res
						return it
					}
				}
			}
			if ops[0] == '+' {
				n := scaffold.Int(ops[2:])
				c.op = func(it item) item {
					res := it.n + n
					it.ops = append(it.ops, fmt.Sprintf("m%v: old(%v) + %v = %v", mid, it.n, n, res))
					it.prevs = append(it.prevs, it.n)
					it.n = res
					return it
				}
			}
		}
		if strings.HasPrefix(l, "Test: divisible by ") {
			n := scaffold.Int(l[len("Test: divisible by "):])
			divprod *= n
			c.test = func(it item) (bool, item) {
				res := it.n%n == 0
				it.ops = append(it.ops, fmt.Sprintf("m%v: %v div by %v: %v", mid, it.n, n, res))
				return res, it
			}
		}
		if strings.HasPrefix(l, "If true: throw to monkey ") {
			n := scaffold.Int(l[len("If true: throw to monkey "):])
			c.nextTrue = n
		}
		if strings.HasPrefix(l, "If false: throw to monkey ") {
			n := scaffold.Int(l[len("If false: throw to monkey "):])
			c.nextFalse = n
		}
	}

	if len(c.items) > 0 {
		monkeys[mid] = c
		mid++
	}

	inspections := make(map[int]int)
	for r := 0; r < 10_000; r++ {
		for i := 0; i < mid; i++ {
			m := monkeys[i]
			for _, it := range m.items {
				inspections[i]++
				w := m.op(it)
				w.n %= divprod
				next := m.nextFalse

				res, it := m.test(w)
				if res {
					next = m.nextTrue
				}
				nm := monkeys[next]
				nm.items = append(nm.items, it)
				monkeys[next] = nm
			}
			m.items = nil
			monkeys[i] = m
		}
		r := r + 1
		if r%1000 == 0 || r == 1 || r == 20 {
			fmt.Printf("%v inspections: %v\n", r, inspections)
		}
	}

	fmt.Printf("inspections: %v\n", inspections)

	vs := maps.Values(inspections)
	sort.Ints(vs)
	bus := vs[len(vs)-1]
	bus *= vs[len(vs)-2]
	fmt.Printf("business: %v\n", bus)
}

type item struct {
	n     int
	prevs []int
	ops   []string
}

type monkey struct {
	items     []item
	op        func(item) item
	test      func(item) (bool, item)
	nextTrue  int
	nextFalse int
}
