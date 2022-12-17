package main

import (
	"fmt"
	"image"
	"strings"

	"github.com/danp/adventofcode/scaffold"
	"golang.org/x/exp/slices"
)

func main() {
	lines := scaffold.Lines()

	x := 1

	var cycle int
	var strength int
	addCycle := func() {
		cycle++
		if slices.Contains([]int{20, 60, 100, 140, 180, 220}, cycle) {
			strength += cycle * x
		}
	}

	for _, l := range lines {
		fs := strings.Fields(l)
		switch fs[0] {
		case "noop":
			addCycle()
		case "addx":
			addCycle()
			addCycle()
			n := scaffold.Int(fs[1])
			x += n
		}
	}

	fmt.Println(strength)

	x = 1
	cycle = 0
	row := -1
	drawn := make(map[image.Point]struct{})
	doCycle := func() {
		pos := cycle % 40
		if pos == 0 {
			row++
		}
		if x == pos || x-1 == pos || x+1 == pos {
			drawn[image.Point{pos, row}] = struct{}{}
		}
		cycle++
	}

	for _, l := range lines {
		fs := strings.Fields(l)
		switch fs[0] {
		case "noop":
			doCycle()
		case "addx":
			doCycle()
			doCycle()
			n := scaffold.Int(fs[1])
			x += n
		}
	}

	for y := 0; y <= row; y++ {
		for x := 0; x < 40; x++ {
			pt := image.Pt(x, y)
			if _, ok := drawn[pt]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
