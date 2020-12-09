package main

import (
	"fmt"
	"sort"
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

	cache := make(map[string]int)
	pq := [][]string{{"shiny gold"}}
	var iters int
	for len(pq) > 0 {
		iters++
		// process longest paths first
		sort.Slice(pq, func(i, j int) bool { return len(pq[i]) > len(pq[j]) })
		path := pq[0]
		k := path[len(path)-1]
		pq = pq[1:]

		down := downgraph[k]
		if len(down) == 0 {
			fmt.Println(path, "is leaf")
			cache[k] = 0
			continue
		}
		allcached := true
		var kbags int
		for o, qty := range down {
			cc, ok := cache[o]
			if ok {
				kbags += qty + (qty * cc)
				continue
			}

			fmt.Println(path, o, "cache miss")
			allcached = false
			newp := make([]string, len(path))
			copy(newp, path)
			newp = append(newp, o)
			pq = append(pq, newp)
		}

		if allcached {
			fmt.Println(path, "now cached with", kbags, "bags")
			cache[k] = kbags
		} else {
			// put ourselves back in the queue so by the time we come
			// back around the items not in the cache are filled in
			pq = append(pq, path)
		}
	}

	fmt.Println(cache["shiny gold"], "bags are required inside shiny gold, iters:", iters)
}
