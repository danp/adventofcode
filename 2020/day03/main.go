package main

import (
	"fmt"
	"image"
	"os"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	var slopes []image.Point

	if os.Args[1] == "part1" {
		slopes = []image.Point{image.Pt(3, 1)}
	} else {
		slopes = []image.Point{image.Pt(1, 1), image.Pt(3, 1), image.Pt(5, 1), image.Pt(7, 1), image.Pt(1, 2)}
	}

	var trees int
	for _, sl := range slopes {
		var sltrees int
		for pt := image.ZP; pt.Y < len(lines); pt = pt.Add(sl) {
			x := pt.X % len(lines[0])
			if lines[pt.Y][x] == '#' {
				sltrees++
			}
		}
		fmt.Println(sl, sltrees)

		if trees == 0 {
			trees = sltrees
		} else {
			trees *= sltrees
		}
	}

	fmt.Println(trees)
}
