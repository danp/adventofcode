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

	scores := make(map[int]int)
	pos := map[int]int{
		1: scaffold.Int(fields0[4]),
		2: scaffold.Int(fields1[4]),
	}

	p := 1
	var rolls int
	for i := 1; i < 1000; i += 3 {
		n := i + (i + 1) + (i + 2)
		np := (pos[p] + n)
		for np > 10 {
			np -= 10
		}
		pos[p] = np
		scores[p] += np

		if scores[p] >= 1000 {
			rolls = i + 2
			break
		}

		switch p {
		case 1:
			p = 2
		case 2:
			p = 1
		}
	}

	fmt.Printf("scores: %v rolls: %v\n", scores, rolls)
}
