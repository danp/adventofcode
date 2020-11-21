package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"io/ioutil"
	"os"

	"github.com/danp/adventofcode/2019/intcode"
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

	if os.Args[1] == "part2" {
		program[0] = 2
	}

	con := controller{}
	gw := newGIFWriter(con.control)

	var g game
	if err := g.run(program, gw.control); err != nil {
		panic(err)
	}

	if os.Args[1] == "part1" {
		var c int
		for _, t := range g.tiles {
			if t == block {
				c++
			}
		}

		fmt.Println(c, "block tiles")
		return
	}

	f, err := os.Create("anim.gif")
	if err != nil {
		panic(err)
	}
	if _, err := gw.WriteTo(f); err != nil {
		panic(err)
	}
	f.Close()
}

type tile int

const (
	empty  tile = 0
	wall   tile = 1
	block  tile = 2
	paddle tile = 3
	ball   tile = 4
)

type game struct {
	vm    *intcode.VM
	tiles map[image.Point]tile

	score int
}

var haltErr = errors.New("halt game")

func (g *game) run(program []int, controller func(g *game) int) error {
	g.tiles = make(map[image.Point]tile)

	input := func() (int, error) {
		co := controller(g)
		if co == -5 {
			return 0, haltErr
		}
		return co, nil
	}

	var (
		ost  int
		x, y int
	)
	output := func(i int) error {
		switch ost {
		case 0:
			x = i
			ost++
		case 1:
			y = i
			ost++
		case 2:
			ost = 0

			if x == -1 && y == 0 {
				fmt.Println("new score", i)
				g.score = i
			} else {
				pt := image.Pt(x, y)
				t := tile(i)
				g.tiles[pt] = t
			}
		}
		return nil
	}

	if g.vm == nil {
		mem := make([]int, 4096)
		g.vm = intcode.NewVM(program, mem, input, output)
	} else {
		g.vm.Input = input
		g.vm.Output = output
	}
	return g.vm.Run()
}

func (g *game) copy() *game {
	ng := &game{
		vm:    g.vm.Copy(),
		tiles: make(map[image.Point]tile),
	}

	for k, v := range g.tiles {
		ng.tiles[k] = v
	}

	return ng
}

func (g *game) findTile(t tile) image.Point {
	for pt, tt := range g.tiles {
		if tt == t {
			return pt
		}
	}
	return image.ZP
}

type controller struct {
	target image.Point
	i      int
}

func (c *controller) control(g *game) int {
	c.i++

	paddle := g.findTile(paddle)

	if c.target.Eq(image.ZP) {
		sc := func(g *game) int {
			ball := g.findTile(ball)
			if ball.Y == paddle.Y-1 {
				c.target = image.Pt(ball.X, paddle.Y)
				return -5 // halt
			}
			return 0
		}

		sg := g.copy()
		// using copied program/etc
		if err := sg.run(nil, sc); err != nil && err != haltErr {
			panic(err)
		}
		fmt.Println(c.i, "new target", c.target)
	} else {
		fmt.Println(c.i, "with target of", c.target, "and paddle at", paddle)
	}

	var out int
	if c.target.X < paddle.X {
		out = -1
	} else if c.target.X > paddle.X {
		out = 1
	}

	ball := g.findTile(ball)
	if ball.Y == paddle.Y-1 {
		fmt.Println(c.i, "ball got to paddle, resetting target")
		c.target = image.ZP
	}

	return out
}

type gifWriter struct {
	inner func(g *game) int

	pal []color.Color
	g   *gif.GIF
}

func newGIFWriter(inner func(g *game) int) *gifWriter {
	return &gifWriter{
		inner: inner,

		g: &gif.GIF{},
		pal: []color.Color{
			color.RGBA{0x00, 0x00, 0x00, 0xff}, color.RGBA{0x00, 0x00, 0xff, 0xff},
			color.RGBA{0x00, 0xff, 0x00, 0xff}, color.RGBA{0x00, 0xff, 0xff, 0xff},
			color.RGBA{0xff, 0x00, 0x00, 0xff}, color.RGBA{0xff, 0x00, 0xff, 0xff},
			color.RGBA{0xff, 0xff, 0x00, 0xff}, color.RGBA{0xff, 0xff, 0xff, 0xff},
		},
	}
}

func (w gifWriter) control(g *game) int {
	img := image.NewPaletted(image.Rect(0, 0, 380, 200), w.pal)

	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			var (
				pt image.Point = image.Pt(x/10, y/10)
				pi int
			)

			switch g.tiles[pt] {
			case wall:
				pi = 1
			case block:
				pi = 2
			case paddle:
				pi = 3
			case ball:
				pi = 4
			}

			img.Set(x, y, w.pal[pi])
		}
	}

	w.g.Image = append(w.g.Image, img)
	w.g.Delay = append(w.g.Delay, 0)

	return w.inner(g)
}

func (w gifWriter) WriteTo(wr io.Writer) (int64, error) {
	var b bytes.Buffer

	if err := gif.EncodeAll(&b, w.g); err != nil {
		return 0, err
	}

	return io.Copy(wr, &b)
}
