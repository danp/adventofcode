package main

import (
	"fmt"
	"image"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	var rect image.Rectangle
	grid := make(map[image.Point]rune)

	var prev image.Point
	for _, l := range lines {
		fs := strings.Split(l, " -> ")
		for i, d := range fs {
			var pt image.Point
			if _, err := fmt.Sscanf(d, "%d,%d", &pt.X, &pt.Y); err != nil {
				panic(err)
			}
			if rect.Min == (image.Point{}) {
				rect.Min = pt
			} else {
				rect.Min.X = min(rect.Min.X, pt.X)
				rect.Min.Y = min(rect.Min.Y, pt.Y)
				rect.Max.X = max(rect.Max.X, pt.X)
				rect.Max.Y = max(rect.Max.Y, pt.Y)
			}

			if i == 0 {
				prev = pt
				continue
			}

			if pt.X == prev.X {
				for y := min(pt.Y, prev.Y); y <= max(pt.Y, prev.Y); y++ {
					grid[image.Point{pt.X, y}] = '#'
				}
			}
			if pt.Y == prev.Y {
				for x := min(pt.X, prev.X); x <= max(pt.X, prev.X); x++ {
					grid[image.Point{x, pt.Y}] = '#'
				}
			}

			prev = pt
		}
	}

	floorY := rect.Max.Y + 2
	show(grid, rect)

	origin := image.Point{500, 0}

outer:
	for {
		sand := origin
	inner:
		for {
			poss := []image.Point{sand.Add(image.Point{0, 1}), sand.Add(image.Point{-1, 1}), sand.Add(image.Point{1, 1})}
			for _, p := range poss {
				if p.Y == floorY {
					continue
				}
				if _, ok := grid[p]; !ok {
					sand = p
					continue inner
				}
			}
			rect.Min.X = min(rect.Min.X, sand.X)
			rect.Min.Y = min(rect.Min.Y, sand.Y)
			rect.Max.X = max(rect.Max.X, sand.X)
			rect.Max.Y = max(rect.Max.Y, sand.Y)

			grid[sand] = 'o'
			if sand == origin {
				break outer
			}
			// show(grid, rect)
			break
		}
	}

	show(grid, rect)

	var resting int
	for _, c := range grid {
		if c == 'o' {
			resting++
		}
	}
	fmt.Printf("resting: %v\n", resting)
}

func show(grid map[image.Point]rune, rect image.Rectangle) {
	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		for x := rect.Min.X; x <= rect.Max.X; x++ {
			pt := image.Point{x, y}
			c := grid[pt]
			if c == 0 {
				c = ' '
			}
			fmt.Print(string(c))
		}
		fmt.Println()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
