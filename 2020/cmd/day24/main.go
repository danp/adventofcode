package main

import (
	"fmt"
	"image"

	"github.com/danp/adventofcode/scaffold"
)

var moves = map[string]image.Point{
	"ne": image.Pt(1, 1),
	"e":  image.Pt(2, 0),
	"se": image.Pt(1, -1),
	"sw": image.Pt(-1, -1),
	"w":  image.Pt(-2, 0),
	"nw": image.Pt(-1, 1),
}

func main() {
	lines := scaffold.Lines()

	grid := make(map[image.Point]bool)

	for _, l := range lines {
		var pt image.Point
		var st int
		for cl := l; len(cl) > 0; {
			for i := 0; i <= len(cl); i++ {
				mvs := cl[:i]
				mv, ok := moves[mvs]
				if !ok {
					continue
				}
				pt = pt.Add(mv)
				st++
				cl = cl[i:]
				break
			}
		}

		// fmt.Println(pt, "from", grid[pt], "to", !grid[pt], "in", st, "steps", l)
		grid[pt] = !grid[pt]
	}

	var black int
	for _, v := range grid {
		if v {
			black++
		}
	}

	// fmt.Println()
	printgrid(grid)
	fmt.Println()
	fmt.Println("initial", black)
	fmt.Println()

	for i := 0; i < 100; i++ {
		newgrid := make(map[image.Point]bool)
		disc := make(map[image.Point]bool)
		var newgridblack int

		process := func(pt image.Point, v bool, dodisc bool) {
			var nb int
			for _, npt := range neighbs(pt) {
				nv, ok := grid[npt]
				if nv {
					nb++
				}
				if dodisc && !ok && !disc[npt] {
					disc[npt] = true
				}
			}

			nv := v
			if v && (nb == 0 || nb > 2) {
				nv = false
			} else if !v && nb == 2 {
				nv = true
			}

			if nv {
				newgridblack++
			}

			newgrid[pt] = nv
		}

		for pt, v := range grid {
			process(pt, v, true)
		}

		for pt := range disc {
			process(pt, false, false)
		}

		grid = newgrid
		// fmt.Println()
		// printgrid(grid)
		// fmt.Println()
		fmt.Println("day", i+1, newgridblack)
		// fmt.Println()
	}
}

func neighbs(pt image.Point) []image.Point {
	out := make([]image.Point, 0, len(moves))
	for _, st := range moves {
		out = append(out, pt.Add(st))
	}
	return out
}

func printgrid(grid map[image.Point]bool) {
	var rect image.Rectangle
	for pt := range grid {
		if !grid[pt] {
			continue
		}
		if pt.X < rect.Min.X {
			rect.Min.X = pt.X
		}
		if pt.Y < rect.Min.Y {
			rect.Min.Y = pt.Y
		}
		if pt.X > rect.Max.X {
			rect.Max.X = pt.X
		}
		if pt.Y > rect.Max.Y {
			rect.Max.Y = pt.Y
		}
	}

	fmt.Printf("%4s", "")
	for x := rect.Min.X; x <= rect.Max.X; x++ {
		x := x
		if x < 0 {
			x *= -1
		}
		fmt.Printf(" %d ", x%10)
	}
	fmt.Println()

	for y := rect.Max.Y; y >= rect.Min.Y; y-- {
		fmt.Printf("%3d ", y)
		for x := rect.Min.X; x <= rect.Max.X; x++ {
			pt := image.Pt(x, y)
			v, ok := grid[pt]
			var pts string
			if !ok {
				pts = " "
			} else {
				if v {
					pts = "b"
				} else {
					pts = "."
				}
			}

			if pt.Eq(image.Point{}) {
				pts = "(" + pts + ")"
			} else {
				pts = " " + pts + " "
			}
			fmt.Print(pts)
		}
		fmt.Println()
	}
}
