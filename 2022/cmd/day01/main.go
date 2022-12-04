package main

import (
	"fmt"
	"sort"

	"github.com/danp/adventofcode/scaffold"
	"golang.org/x/exp/maps"
)

func main() {
	lines := scaffold.Lines()

	elfCals := make(map[int]int)
	var n int
	for _, l := range lines {
		if l == "" {
			n++
			continue
		}
		elfCals[n] += scaffold.Int(l)
	}

	vals := maps.Values(elfCals)
	sort.Sort(sort.Reverse(sort.IntSlice(vals)))

	fmt.Println(vals[0])

	var top3 int
	for _, v := range vals[0:3] {
		top3 += v
	}
	fmt.Println(top3)
}
