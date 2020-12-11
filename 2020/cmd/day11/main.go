package main

import (
	"fmt"
	"image"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	origGrid := make(map[image.Point]bool)
	var max image.Point
	for y, l := range lines {
		for x, c := range l {
			pt := image.Pt(x, y)
			max = pt
			if c == 'L' {
				origGrid[pt] = false
			}
		}
	}

	slopes := []image.Point{
		image.Pt(1, 0),   // right
		image.Pt(-1, 0),  // left
		image.Pt(0, 1),   // down
		image.Pt(0, -1),  // up
		image.Pt(-1, -1), // up-left
		image.Pt(1, -1),  // up-right
		image.Pt(-1, 1),  // down-left
		image.Pt(1, 1),   // down-right
	}

	grid := copyGrid(origGrid)
	for {
		newgrid := make(map[image.Point]bool)
		var changes int
		for pt, ptocc := range grid {
			var neighbs int

			for _, sl := range slopes {
				if grid[pt.Add(sl)] {
					neighbs++
				}
			}

			if !ptocc && neighbs == 0 {
				newgrid[pt] = true
			} else if ptocc && neighbs >= 4 {
				newgrid[pt] = false
			} else {
				newgrid[pt] = grid[pt]
			}

			if newgrid[pt] != grid[pt] {
				changes++
			}
		}
		if changes == 0 {
			break
		}
		grid = newgrid
	}

	var occ int
	for _, v := range grid {
		if v {
			occ++
		}
	}
	fmt.Println("seats occupied with initial rules:", occ)

	grid = copyGrid(origGrid)

	// grid point -> slope -> point of first seat visible along that slope
	ptslopecache := make(map[image.Point]map[image.Point]image.Point)
	for pt := range grid {
		ptslopecache[pt] = make(map[image.Point]image.Point)
	}

	for {
		newgrid := make(map[image.Point]bool)
		var changes int
		for pt, ptocc := range grid {
			var neighbs int

			for _, sl := range slopes {
				slpt, ok := ptslopecache[pt][sl]
				if ok {
					if grid[slpt] {
						neighbs++
					}
					continue
				}

				for slpt := pt.Add(sl); slpt.X >= 0 && slpt.X <= max.X && slpt.Y >= 0 && slpt.Y <= max.Y; slpt = slpt.Add(sl) {
					v, ok := grid[slpt]
					if v {
						neighbs++
					}
					if ok {
						ptslopecache[pt][sl] = slpt
						break
					}
				}
			}

			if !ptocc && neighbs == 0 {
				newgrid[pt] = true
			} else if ptocc && neighbs >= 5 {
				newgrid[pt] = false
			} else {
				newgrid[pt] = grid[pt]
			}

			if newgrid[pt] != grid[pt] {
				changes++
			}
		}
		if changes == 0 {
			break
		}
		grid = newgrid
	}

	occ = 0
	for _, v := range grid {
		if v {
			occ++
		}
	}
	fmt.Println("seats occupied with refined rules:", occ)
}

func display(grid map[image.Point]bool, max image.Point) {
	for y := 0; y <= max.Y; y++ {
		for x := 0; x <= max.X; x++ {
			occ, ok := grid[image.Pt(x, y)]
			if !ok {
				fmt.Print(".")
				continue
			}
			switch occ {
			case true:
				fmt.Print("#")
			case false:
				fmt.Print("L")
			}
		}
		fmt.Println()
	}
}

func copyGrid(grid map[image.Point]bool) map[image.Point]bool {
	newgrid := make(map[image.Point]bool)
	for k, v := range grid {
		newgrid[k] = v
	}
	return newgrid
}
