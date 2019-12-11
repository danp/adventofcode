package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/danp/adventofcode2019/intcode"
)

const (
	up    = 0
	right = 1
	down  = 2
	left  = 3

	black = 0
	white = 1
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	program, err := intcode.Parse(string(b))
	if err != nil {
		panic(err)
	}

	mem := make([]int, 2048)

	var (
		loc image.Point
		dir int = up
		ost int
	)

	painted := make(map[image.Point]int)
	if os.Args[1] == "part2" {
		painted[loc] = white
	}

	input := func() (int, error) {
		return painted[loc], nil
	}

	output := func(x int) error {
		if ost == 0 {
			painted[loc] = x
			ost++
			return nil
		}
		ost = 0

		switch x {
		case 0:
			dir--
		case 1:
			dir++
		}
		switch dir {
		case -1:
			dir = left
		case 4:
			dir = up
		}

		switch dir {
		case up:
			loc.Y--
		case right:
			loc.X++
		case down:
			loc.Y++
		case left:
			loc.X--
		}

		return nil
	}

	if err := intcode.Run(program, mem, input, output); err != nil {
		panic(err)
	}

	if os.Args[1] == "part2" {
		img := mapImage(painted)
		if err := png.Encode(os.Stdout, img); err != nil {
			panic(err)
		}
	} else {
		fmt.Println(len(painted))
	}
}

type mapImage map[image.Point]int

func (l mapImage) ColorModel() color.Model {
	return color.GrayModel
}

func (l mapImage) Bounds() image.Rectangle {
	var min, max image.Point

	for p := range l {
		if p.X < min.X {
			min.X = p.X
		}
		if p.Y < min.Y {
			min.Y = p.Y
		}
		if p.X > max.X {
			max.X = p.X
		}
		if p.Y > max.Y {
			max.Y = p.Y
		}
	}

	return image.Rect(min.X, min.Y, max.X+1, max.Y+1)
}

func (l mapImage) At(x, y int) color.Color {
	pt := image.Pt(x, y)
	return color.Gray{Y: uint8(l[pt]) * 255}
}
