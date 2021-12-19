package main

import (
	"fmt"
	"image"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	grid := make(map[image.Point]int)
	for y, l := range lines {
		for x, c := range l {
			grid[image.Pt(x, y)] = int(c) - '0'
		}
	}

	var totalFlashes int
	stepGrid := grid
	for i := 0; i < 100; i++ {
		newGrid, flashes := step(stepGrid)
		totalFlashes += flashes
		stepGrid = newGrid
	}
	fmt.Printf("totalFlashes: %v\n", totalFlashes)

	stepGrid = grid
	for i := 0; i < 10000; i++ {
		newGrid, flashes := step(stepGrid)
		if flashes == len(newGrid) {
			fmt.Printf("sync flash step: %v\n", i+1)
			break
		}
		stepGrid = newGrid
	}
}

func step(grid map[image.Point]int) (map[image.Point]int, int) {
	newGrid := make(map[image.Point]int)
	flashers := make(map[image.Point]struct{})
	todo := make(map[image.Point]struct{})
	for pt, l := range grid {
		nl := l + 1
		newGrid[pt] = nl
		if nl > 9 {
			flashers[pt] = struct{}{}
			todo[pt] = struct{}{}
		}
	}

	for len(todo) > 0 {
		for pt := range todo {
			delete(todo, pt)

			for _, npt := range neigbs(pt) {
				l, ok := newGrid[npt]
				if !ok {
					continue
				}
				nl := l + 1
				if nl > 9 {
					if _, ok := flashers[npt]; !ok {
						flashers[npt] = struct{}{}
						todo[npt] = struct{}{}
					}
				}
				newGrid[npt] = nl
			}
		}
	}

	for pt := range flashers {
		newGrid[pt] = 0
	}

	return newGrid, len(flashers)
}

func neigbs(pt image.Point) []image.Point {
	dirs := []image.Point{
		{0, -1},  // n
		{1, -1},  // ne
		{1, 0},   // e
		{1, 1},   // se
		{0, 1},   // s
		{-1, 1},  // sw
		{-1, 0},  // w
		{-1, -1}, // nw
	}

	var pts []image.Point
	for _, d := range dirs {
		pts = append(pts, pt.Add(d))
	}
	return pts
}
