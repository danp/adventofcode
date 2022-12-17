package main

import (
	"fmt"
	"image"
	"math"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	do(lines, []image.Point{{}, {}})

	var rope []image.Point
	for i := 0; i < 10; i++ {
		rope = append(rope, image.Point{})
	}
	do(lines, rope)
}

func do(lines []string, rope []image.Point) {
	touched := make(map[image.Point]struct{})

	for _, l := range lines {
		for i := 0; i < scaffold.Int(l[2:]); i++ {
			var add image.Point
			switch l[0:1] {
			case "L":
				add = image.Pt(-1, 0)
			case "R":
				add = image.Pt(1, 0)
			case "U":
				add = image.Pt(0, -1)
			case "D":
				add = image.Pt(0, 1)
			}
			rope[0] = rope[0].Add(add)

			for j := 1; j < len(rope); j++ {
				prev := rope[j-1]
				curr := rope[j]
				diff := image.Point{prev.X - curr.X, prev.Y - curr.Y}

				if abs(diff.X) <= 1 && abs(diff.Y) <= 1 {
					//
				} else {
					// head: (2, 0) tail: (0, 0)
					// diff: (2, 0)
					// want: tail: (1, 0)

					// head: (2, 0) tail: (4, 0)
					// diff: (-2, 0)
					// want: tail: (3, 0)

					// head: (4, 0) tail: (3, 2)
					// diff: (1, -2)
					// want: tail: (4, 1)
					// want add: (1, -1)

					if diff.X < -1 {
						diff.X += 1
					} else if diff.X > 1 {
						diff.X -= 1
					}
					if diff.Y < -1 {
						diff.Y += 1
					} else if diff.Y > 1 {
						diff.Y -= 1
					}

					rope[j].Y += diff.Y
					rope[j].X += diff.X
				}
			}

			touched[rope[len(rope)-1]] = struct{}{}
		}
	}

	fmt.Printf("len(touched): %v\n", len(touched))
}

func abs(x int) int {
	return int(math.Abs(float64(x)))
}
