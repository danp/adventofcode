package main

import (
	"fmt"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	min := 1
	var max int

	cupss := lines[0]
	cups := make(map[int]int)
	var prev, first int
	for i, c := range cupss {
		ci := scaffold.Int(string(c))
		if ci > max {
			max = ci
		}
		if i == 0 {
			first = ci
		} else {
			cups[prev] = ci
		}
		prev = ci
	}
	cups[prev] = first

	inputCups := make(map[int]int)
	for k, v := range cups {
		inputCups[k] = v
	}

	play(inputCups, first, 100, min, max)
	show("input-100", inputCups)

	bigCups := make(map[int]int)
	for k, v := range cups {
		bigCups[k] = v
	}
	for i := max + 1; i <= 1_000_000; i++ {
		bigCups[prev] = i
		prev = i
	}
	bigCups[prev] = first

	play(bigCups, first, 10000000, min, 1_000_000)
	fmt.Println(bigCups[1], bigCups[bigCups[1]], bigCups[1]*bigCups[bigCups[1]])
}

func play(cups map[int]int, cur, moves, min, max int) {
	for move := 0; move < moves; move++ {
		fpu := cups[cur]
		spu := cups[fpu]
		tpu := cups[spu]

		dst := cur - 1
		if dst < min {
			dst = max
		}

		for dst == fpu || dst == spu || dst == tpu {
			dst--
			if dst < min {
				dst = max
			}
		}

		apu := cups[tpu]

		cups[cur] = apu

		ad := cups[dst]
		cups[dst] = fpu
		cups[tpu] = ad

		cur = cups[cur]
	}
}

func show(label string, cups map[int]int) {
	fmt.Print(label, " ")
	k := 1
	for {
		v := cups[k]
		if v == 1 {
			break
		}
		fmt.Print(v)
		k = v
	}
	fmt.Println()
}
