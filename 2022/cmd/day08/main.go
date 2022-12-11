package main

import (
	"fmt"
	"image"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	var rect image.Rectangle
	grid := make(map[image.Point]int)
	for y, l := range lines {
		for x, c := range l {
			pt := image.Point{x, y}
			grid[pt] = int(c - '0')
			rect.Max = pt
		}
	}

	visible := func(pt image.Point) bool {
	X1:
		for x := rect.Min.X; x <= pt.X; x++ {
			tpt := image.Point{x, pt.Y}
			if tpt == pt {
				return true
			}
			if grid[tpt] >= grid[pt] {
				break X1
			}
		}
	X2:
		for x := rect.Max.X; x >= pt.X; x-- {
			tpt := image.Point{x, pt.Y}
			if tpt == pt {
				return true
			}
			if grid[tpt] >= grid[pt] {
				break X2
			}
		}
	Y1:
		for y := rect.Min.Y; y <= pt.Y; y++ {
			tpt := image.Point{pt.X, y}
			if tpt == pt {
				return true
			}
			if grid[tpt] >= grid[pt] {
				break Y1
			}
		}
		for y := rect.Max.Y; y >= pt.Y; y-- {
			tpt := image.Point{pt.X, y}
			if tpt == pt {
				return true
			}
			if grid[tpt] >= grid[pt] {
				return false
			}
		}
		return true
	}

	vis := make(map[image.Point]struct{})
	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		for x := rect.Min.X; x <= rect.Max.X; x++ {
			pt := image.Point{x, y}

			if visible(pt) {
				vis[pt] = struct{}{}
			}
		}
	}

	fmt.Printf("len(vis): %v\n", len(vis))

	score := func(pt image.Point) int {
		var s int
	X1:
		for x := pt.X; x <= rect.Max.X; x++ {
			tpt := image.Point{x, pt.Y}
			if tpt == pt {
				continue
			}
			s++
			if grid[tpt] >= grid[pt] {
				break X1
			}
		}
		var n int
	X2:
		for x := pt.X; x >= rect.Min.X; x-- {
			tpt := image.Point{x, pt.Y}
			if tpt == pt {
				continue
			}
			n++
			if grid[tpt] >= grid[pt] {
				break X2
			}
		}
		s *= n
		n = 0
	Y1:
		for y := pt.Y; y <= rect.Max.Y; y++ {
			tpt := image.Point{pt.X, y}
			if tpt == pt {
				continue
			}
			n++
			if grid[tpt] >= grid[pt] {
				break Y1
			}
		}
		s *= n
		n = 0
	Y2:
		for y := pt.Y; y >= rect.Min.Y; y-- {
			tpt := image.Point{pt.X, y}
			if tpt == pt {
				continue
			}
			n++
			if grid[tpt] >= grid[pt] {
				break Y2
			}
		}
		return s * n
	}

	var max int
	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		for x := rect.Min.X; x <= rect.Max.X; x++ {
			pt := image.Point{x, y}

			s := score(pt)
			if s > max {
				max = s
			}
		}
	}

	fmt.Printf("max: %v\n", max)
}
