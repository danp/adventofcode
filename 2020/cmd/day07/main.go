package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	upgraph := make(map[string]map[string]int)
	downgraph := make(map[string]map[string]int)

	for _, l := range lines {
		parts := strings.SplitN(l, " bags contain ", 2)

		name := parts[0]
		if upgraph[name] == nil {
			upgraph[name] = make(map[string]int)
		}
		if downgraph[name] == nil {
			downgraph[name] = make(map[string]int)
		}
		if parts[1] == "no other bags." {
			continue
		}
		for _, p := range strings.Split(parts[1], ", ") {
			f := strings.Fields(p)
			qty, err := strconv.Atoi(f[0])
			if err != nil {
				panic(err)
			}
			cont := strings.Join(f[1:len(f)-1], " ")

			if upgraph[cont] == nil {
				upgraph[cont] = make(map[string]int)
			}
			upgraph[cont][name] = qty
			downgraph[name][cont] = qty
		}
	}

	q := []string{"shiny gold"}
	seen := make(map[string]bool)
	for len(q) > 0 {
		k := q[0]
		q = q[1:]

		up := upgraph[k]
		if len(up) == 0 {
			continue
		}

		for o := range up {
			if !seen[o] {
				seen[o] = true
				q = append(q, o)
			}
		}
	}

	fmt.Println(len(seen), "bag colors can eventually contain shiny gold")

	q = []string{"shiny gold"}
	var bags int
	for len(q) > 0 {
		k := q[0]
		q = q[1:]

		down := downgraph[k]
		if len(down) == 0 {
			continue
		}
		for o, qty := range down {
			for i := 0; i < qty; i++ {
				q = append(q, o)
			}
			bags += qty
		}
	}

	fmt.Println(bags, "bags are required inside shiny gold")
}
