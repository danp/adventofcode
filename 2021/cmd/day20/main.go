package main

import (
	"fmt"
	"image"
	"math"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	var algo string
	for lines[0] != "" {
		algo += lines[0]
		lines = lines[1:]
	}
	fmt.Printf("algo: %d %v\n", len(algo), algo)
	fmt.Println()

	var rect image.Rectangle
	grid := make(map[image.Point]bool)
	for y, l := range lines[1:] {
		for x, c := range l {
			pt := image.Pt(x, y)
			if c == '#' {
				grid[pt] = true
			}
			rect.Max = pt
		}
	}

	// first input is default=dark
	// pixels outside boundary will flip to whatever algo[0] is
	// pixels inside boundary will assume missing are dark, change based on that
	//
	// result: missing=lit, known pixels are what they are
	//
	// this input is default=lit
	// pixels outside boundary will flip to whatever algo[-1] is
	// pixels inside boundary will assume missing are lit, change based on that

	show(rect, grid, false)

	var missing bool
	for i := 0; i < 2; i++ {
		grid, rect, missing = step(algo, rect, grid, missing)
		show(rect, grid, missing)
	}

	var lit int
	for _, v := range grid {
		if v {
			lit++
		}
	}
	fmt.Printf("lit: %v\n", lit)
}

func show(rect image.Rectangle, grid map[image.Point]bool, missing bool) {
	for y := rect.Min.Y - 5; y <= rect.Max.Y+5; y++ {
		for x := rect.Min.X - 5; x <= rect.Max.X+5; x++ {
			if v, ok := grid[image.Pt(x, y)]; v || (!ok && missing) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func step(algo string, rect image.Rectangle, grid map[image.Point]bool, missing bool) (map[image.Point]bool, image.Rectangle, bool) {
	newGrid := make(map[image.Point]bool)
	rect.Min = rect.Min.Add(image.Pt(-1, -1))
	rect.Max = rect.Max.Add(image.Pt(1, 1))
	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		for x := rect.Min.X; x <= rect.Max.X; x++ {
			pt := image.Pt(x, y)
			i := neigbsToInt(grid, pt, missing)
			newGrid[pt] = algo[i] == '#'
		}
	}

	if missing {
		missing = algo[len(algo)-1] == '#'
	} else {
		missing = algo[0] == '#'
	}

	return newGrid, rect, missing
}

func neigbsToInt(grid map[image.Point]bool, pt image.Point, missing bool) int {
	var bits []int
	for _, npt := range neigbs(pt) {
		var b int
		if v, ok := grid[npt]; v || (!ok && missing) {
			b = 1
		}
		bits = append(bits, b)
	}
	n := bitsToInt(bits)
	return n
}

func neigbs(pt image.Point) []image.Point {
	var out []image.Point
	for y := pt.Y - 1; y <= pt.Y+1; y++ {
		for x := pt.X - 1; x <= pt.X+1; x++ {
			out = append(out, image.Pt(x, y))
		}
	}
	return out
}

func bitsToInt(bits []int) int {
	var n int
	for i := 0; i < len(bits); i++ {
		b := bits[len(bits)-1-i]
		if b == 1 {
			n += int(math.Pow(2, float64(i)))
		}
	}
	return n
}
