package main

import (
	"fmt"
	"image"
	"math"
	"strconv"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	commands := make([]command, 0, len(lines))
	for _, l := range lines {
		ins := l[0]
		arg, err := strconv.Atoi(l[1:])
		if err != nil {
			panic(err)
		}
		commands = append(commands, command{ins: rune(ins), arg: arg})
	}

	headings := []image.Point{
		image.Pt(1, 0),  // east
		image.Pt(0, 1),  // south
		image.Pt(-1, 0), // west
		image.Pt(0, -1), // north
	}

	{
		var pt image.Point
		headIdx := 0
		for _, c := range commands {
			switch c.ins {
			case 'N':
				pt = pt.Add(image.Pt(0, -c.arg))
			case 'S':
				pt = pt.Add(image.Pt(0, c.arg))
			case 'E':
				pt = pt.Add(image.Pt(c.arg, 0))
			case 'W':
				pt = pt.Add(image.Pt(-c.arg, 0))
			case 'F':
				pt = pt.Add(headings[headIdx].Mul(c.arg))
			case 'L', 'R':
				rot := c.arg / 90
				if c.ins == 'L' {
					rot *= -1
				}
				headIdx += rot
				if headIdx >= len(headings) {
					headIdx -= len(headings)
				}
				if headIdx < 0 {
					headIdx += len(headings)
				}
			}
		}

		fmt.Println(pt, math.Abs(float64(pt.X))+math.Abs(float64(pt.Y)))
	}
	fmt.Println()

	var pt image.Point
	waypoint := pt.Add(image.Pt(10, 1))

	for _, c := range commands {
		switch c.ins {
		case 'N':
			waypoint = waypoint.Add(image.Pt(0, c.arg))
		case 'S':
			waypoint = waypoint.Add(image.Pt(0, -c.arg))
		case 'E':
			waypoint = waypoint.Add(image.Pt(c.arg, 0))
		case 'W':
			waypoint = waypoint.Add(image.Pt(-c.arg, 0))
		case 'F':
			shift := waypoint.Sub(pt).Mul(c.arg)
			pt = pt.Add(shift)
			waypoint = waypoint.Add(shift)
		case 'R', 'L':
			diff := pt.Sub(waypoint)
			if c.ins == 'L' {
				c.ins = 'R'
				c.arg = 360 - c.arg
			}
			newdiff := diff
			for i := 0; i < c.arg/90; i++ {
				newdiff = image.Pt(-newdiff.Y, newdiff.X)
			}
			if c.arg == 180 {
				newdiff = newdiff.Mul(-1)
			}
			waypoint = pt.Add(newdiff)
		}
	}

	fmt.Println(pt, math.Abs(float64(pt.X))+math.Abs(float64(pt.Y)))
}

type command struct {
	ins rune
	arg int
}
