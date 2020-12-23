package main

import (
	"fmt"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	cupss := lines[0]
	var cups []int
	for _, c := range cupss {
		cups = append(cups, scaffold.Int(string(c)))
	}

	clock1 := make(map[int]int)
	var seq1 []int

	const min, max = 1, 9
	var cur int
	for turn := 0; turn < 100; turn++ {
		pu := make([]int, 3)
		for i := 0; i < 3; i++ {
			pi := i + cur + 1
			if pi >= len(cups) {
				pi -= len(cups)
			}
			pu[i] = cups[pi]
		}

		dst := cups[cur] - 1
		if dst < min {
			dst = max
		}
		for contains(pu, dst) {
			dst--
			if dst < min {
				dst = max
			}
		}
		di := index(cups, dst)
		newdi := di - len(pu)
		if newdi < 0 {
			newdi += len(cups)
		}
		// fmt.Println(turn, cups, cur, cups[cur], pu, dst, newdi)

		newcups := make([]int, len(cups))
		newcups[newdi] = dst
		wrotepu := make([]int, 3)
		for i, puu := range pu {
			ni := newdi + 1 + i
			if ni >= len(cups) {
				ni -= len(cups)
			}
			newcups[ni] = puu
			wrotepu[i] = ni
		}

		for i := 0; i < len(cups); i++ {
			if i == newdi || contains(wrotepu, i) {
				continue
			}

			oi := i
			if (cur < newdi && i > cur && i < newdi) || (cur > newdi && (i > cur || i < newdi)) {
				oi += 3
				if oi >= len(cups) {
					oi -= len(cups)
				}
			}
			newcups[i] = cups[oi]
		}

		ncc := make(map[int]bool)
		for _, v := range newcups {
			if ncc[v] {
				fmt.Println(newcups)
				panic(fmt.Sprintf("dup value %d", v))
			}
			ncc[v] = true
		}

		cur++
		if cur >= len(cups) {
			cur = 0
		}
		// cups, newcups = newcups, cups
		cups = newcups

		idx1 := index(cups, 1) + 1
		if idx1 >= len(cups) {
			idx1 -= len(cups)
		}
		clock1[cups[idx1]]++
		seq1 = append(seq1, cups[idx1])
		fmt.Println(turn, "1c", cups[idx1], clock1)
	}

	fmt.Println(cups)
	fmt.Println(clock1)
	fmt.Println(seq1)

	var out string
	for i := 0; i < len(cups)-1; i++ {
		out += fmt.Sprint(cups[(1+i+index(cups, 1))%len(cups)])
	}
	fmt.Println(out)
}

func contains(s []int, x int) bool {
	for _, n := range s {
		if n == x {
			return true
		}
	}
	return false
}

func index(s []int, x int) int {
	for i, n := range s {
		if n == x {
			return i
		}
	}
	return -1
}
