package main

import (
	"fmt"
	"image"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	grid := make(map[image.Point]struct{})
	var inFolds bool
	var folds []image.Point
	for _, l := range lines {
		if l == "" {
			inFolds = true
			continue
		}

		if inFolds {
			parts := strings.Fields(l)
			a, b, _ := strings.Cut(parts[2], "=")
			n := scaffold.Int(b)
			var pt image.Point
			switch a {
			case "x":
				pt.X = n
			default:
				pt.Y = n
			}
			folds = append(folds, pt)
			continue
		}

		parts := strings.Split(l, ",")
		pt := image.Pt(scaffold.Int(parts[0]), scaffold.Int(parts[1]))
		grid[pt] = struct{}{}
	}

	for i, f := range folds {
		grid = fold(grid, f)
		fmt.Println("after fold", i+1, "there are", len(grid), "image.Points visible")
	}

	r := rect(grid)
	for y := r.Min.Y; y <= r.Max.Y; y++ {
		for x := r.Min.X; x <= r.Max.X; x++ {
			if _, ok := grid[image.Pt(x, y)]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func fold(grid map[image.Point]struct{}, line image.Point) map[image.Point]struct{} {
	translate := func(pt image.Point) image.Point {
		if line.X > 0 && pt.X > line.X {
			return pt.Add(image.Pt(-(pt.X-line.X)*2, 0))
		} else if line.Y > 0 && pt.Y > line.Y {
			return pt.Add(image.Pt(0, -(pt.Y-line.Y)*2))
		} else {
			return pt
		}
	}

	newGrid := make(map[image.Point]struct{})
	for pt := range grid {
		newGrid[translate(pt)] = struct{}{}
	}
	return newGrid
}

func rect(grid map[image.Point]struct{}) image.Rectangle {
	var r image.Rectangle
	for pt := range grid {
		r.Max.X = max(r.Max.X, pt.X)
		r.Max.Y = max(r.Max.Y, pt.Y)
	}
	return r
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
