package main

import (
	"fmt"
	"image"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	var target image.Rectangle
	targetFields := strings.Fields(lines[0])
	_, xs, _ := strings.Cut(targetFields[2][:len(targetFields[2])-1], "=")
	x1, x2, _ := strings.Cut(xs, "..")
	target.Min.X = scaffold.Int(x1)
	target.Max.X = scaffold.Int(x2)
	_, ys, _ := strings.Cut(targetFields[3], "=")
	y1, y2, _ := strings.Cut(ys, "..")
	target.Min.Y = scaffold.Int(y1)
	target.Max.Y = scaffold.Int(y2)
	fmt.Printf("target: %v\n", target)

	start := image.Point{}

	var hits []image.Point
	var bestY int
	for x := 0; x < 10000; x++ {
		for y := -100; y < 1000; y++ {
			vel := image.Pt(x, y)

			maxY, _, ok := check(start, vel, target)
			if ok {
				hits = append(hits, vel)
				if maxY > bestY {
					bestY = maxY
					fmt.Println("new best", bestY)
				}
			}
		}
	}

	fmt.Println(len(hits))
}

func check(start, vel image.Point, target image.Rectangle) (int, image.Point, bool) {
	var maxY int
	pos := start

	step := func() {
		pos.X += vel.X
		pos.Y += vel.Y

		if pos.Y > maxY {
			maxY = pos.Y
		}

		if vel.X > 0 {
			vel.X -= 1
		}
		vel.Y -= 1
	}

	for i := 0; i < 100000; i++ {
		step()
		if pos.In(target) {
			return maxY, pos, true
		}

		if pos.Y < target.Min.Y-10 {
			return maxY, pos, false
		}

		if pos.X > target.Max.X+100 {
			return maxY, pos, false
		}
	}

	fmt.Printf("target: %v pos: %v\n", target, pos)

	panic("broke")
}
