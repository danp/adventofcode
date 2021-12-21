package main

import (
	"fmt"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()
	fields0 := strings.Fields(lines[0])
	fields1 := strings.Fields(lines[1])

	pos := map[int]int{
		1: scaffold.Int(fields0[4]),
		2: scaffold.Int(fields1[4]),
	}

	scores, wp, rolls := play(pos, 1000, func(i int) int {
		return i + (i + 1) + (i + 2)
	})
	fmt.Printf("scores: %v wp: %v rolls: %v\n", scores, wp, rolls)

	for i := 3; i <= 9; i++ {
		scores, wp, rolls = play(pos, 21, func(int) int {
			return i
		})
		fmt.Printf("scores: %v wp: %v rolls: %v\n", scores, wp, rolls)
	}
}

func play(pos map[int]int, winCutoff int, roll func(int) int) (map[int]int, int, int) {
	npos := make(map[int]int)
	for k, v := range pos {
		npos[k] = v
	}
	pos = npos
	scores := make(map[int]int)
	p := 1
	var rolls int
	for i := 1; i < 1000; i += 3 {
		n := roll(i)
		np := (pos[p] + n)
		for np > 10 {
			np -= 10
		}
		pos[p] = np
		scores[p] += np

		if scores[p] >= winCutoff {
			rolls = i + 2
			return scores, p, rolls
		}

		switch p {
		case 1:
			p = 2
		case 2:
			p = 1
		}
	}
	panic("no winner")
}
