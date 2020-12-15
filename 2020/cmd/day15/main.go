package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	hist := make(map[int][]int)
	var turn, prev int

	for _, ns := range strings.Split(lines[0], ",") {
		turn++
		n, _ := strconv.Atoi(ns)
		hist[n] = append(hist[n], turn)
		fmt.Println(turn, n)
		prev = n
	}

	for {
		turn++

		ph := hist[prev]

		var n int // 0 by default
		if len(ph) > 1 {
			n = ph[len(ph)-1] - ph[len(ph)-2]
		}
		hist[n] = append(hist[n], turn)

		prev = n

		if turn == 2020 {
			fmt.Println(turn, n)
		}
		if turn == 30000000 {
			fmt.Println(turn, n)
			break
		}
	}
}
