package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/danp/adventofcode2019/intcode"
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

	if os.Args[1] == "part1" {
		control := newControl()

		mem := make([]int, 4096)
		if err := intcode.Run(program, mem, control.input, control.output); err != nil {
			panic(err)
		}

		if control.oxy.Eq(image.ZP) {
			panic("didn't find oxygen")
		}
		trace := solve(control)
		ng := copyGrid(control.grid)
		for _, p := range trace {
			ng[p] = path
		}
		display(ng, control.pos, control.dir)
		fmt.Println(len(trace) - 1) // starts where we are
	} else {
		control := newExpansiveControl()
		mem := make([]int, 4096)
		if err := intcode.Run(program, mem, control.input, control.output); err != nil {
			panic(err)
		}

		grid := copyGrid(control.grid)
		fmt.Println(expandOxygen(grid))
		display(grid, image.ZP, north)
	}
}

type content int

const (
	unknown content = 0
	travel  content = 1
	wall    content = 2
	oxygen  content = 3

	path content = 10
)

type dir int

const (
	north dir = 1
	south dir = 2
	west  dir = 3
	east  dir = 4
)

type control struct {
	pos  image.Point
	dir  dir
	grid map[image.Point]content
	oxy  image.Point
}

func newControl() *control {
	return &control{
		dir: north,
		grid: map[image.Point]content{
			image.ZP: travel,
		},
	}
}

func (c *control) input() (int, error) {
	if !c.oxy.Eq(image.ZP) {
		display(c.grid, c.pos, c.dir)
		// halt
		return 0, nil
	}
	return int(c.dir), nil
}

func (c *control) output(x int) error {
	var cand image.Point

	switch c.dir {
	case north:
		cand = c.pos.Add(image.Pt(0, -1))
	case south:
		cand = c.pos.Add(image.Pt(0, 1))
	case west:
		cand = c.pos.Add(image.Pt(-1, 0))
	case east:
		cand = c.pos.Add(image.Pt(1, 0))
	}

	switch x {
	case 0:
		c.grid[cand] = wall
		c.turn()
	case 1:
		c.grid[cand] = travel
		c.pos = cand
		c.maybeTurn()
	case 2:
		c.grid[cand] = oxygen
		c.pos = cand
		c.oxy = cand
	}

	return nil
}

func (c *control) turn() {
	c.dir += dir(1 + rand.Intn(2))
	if c.dir > east {
		c.dir = c.dir - east
	}
}

func (c *control) maybeTurn() {
	if rand.Intn(10) <= 2 {
		c.turn()
	}
}

type expansiveControl struct {
	pos  image.Point
	dir  dir
	grid map[image.Point]content

	calls int
}

func newExpansiveControl() *expansiveControl {
	return &expansiveControl{
		dir: north,
		grid: map[image.Point]content{
			image.ZP: travel,
		},
	}
}

func (e *expansiveControl) input() (int, error) {
	up := findPath(
		e.grid,
		e.pos,
		func(path []image.Point) bool {
			con := e.grid[path[len(path)-1]]
			return con == unknown
		},
		func(pt image.Point) bool {
			con := e.grid[pt]
			return con == travel || con == unknown
		},
	)

	if len(up) < 2 {
		display(e.grid, e.pos, e.dir)
		return 0, nil
	}

	np := up[1]
	var d dir
	if np.Y < e.pos.Y {
		d = north
	} else if np.Y > e.pos.Y {
		d = south
	} else if np.X < e.pos.X {
		d = west
	} else if np.X > e.pos.X {
		d = east
	}
	e.dir = d

	return int(d), nil
}

func (e *expansiveControl) output(x int) error {
	var cand image.Point

	switch e.dir {
	case north:
		cand = e.pos.Add(image.Pt(0, -1))
	case south:
		cand = e.pos.Add(image.Pt(0, 1))
	case west:
		cand = e.pos.Add(image.Pt(-1, 0))
	case east:
		cand = e.pos.Add(image.Pt(1, 0))
	}

	switch x {
	case 0:
		e.grid[cand] = wall
	case 1:
		e.grid[cand] = travel
		e.pos = cand
	case 2:
		e.grid[cand] = oxygen
		e.pos = cand
	}

	return nil
}

func solve(c *control) []image.Point {
	return findPath(
		c.grid,
		c.oxy,
		func(path []image.Point) bool {
			return path[len(path)-1].Eq(image.ZP)
		},
		func(pt image.Point) bool {
			return c.grid[pt] == travel
		},
	)
}

// findPath finds the shortest path from start to where check returns true.
func findPath(grid map[image.Point]content, start image.Point, check func(path []image.Point) bool, consider func(pt image.Point) bool) []image.Point {
	seen := make(map[image.Point]bool)

	q := [][]image.Point{{start}}
	for len(q) > 0 {
		path := q[0]
		q = q[1:]

		if check(path) {
			return path
		}

		pt := path[len(path)-1]

		for _, cpt := range []image.Point{
			pt.Add(image.Pt(0, -1)), // north
			pt.Add(image.Pt(0, 1)),  // south
			pt.Add(image.Pt(-1, 0)), // west
			pt.Add(image.Pt(1, 0)),  // east
		} {
			if !consider(cpt) || seen[cpt] {
				continue
			}

			seen[cpt] = true

			newp := make([]image.Point, len(path))
			copy(newp, path)
			newp = append(newp, cpt)
			q = append(q, newp)
		}

	}

	return nil
}

func copyGrid(grid map[image.Point]content) map[image.Point]content {
	ng := make(map[image.Point]content)
	for k, v := range grid {
		ng[k] = v
	}
	return ng
}

func display(grid map[image.Point]content, droidPos image.Point, droidDir dir) {
	var min, max image.Point

	for pt := range grid {
		if pt.X < min.X {
			min.X = pt.X
		}
		if pt.Y < min.Y {
			min.Y = pt.Y
		}
		if pt.X > max.X {
			max.X = pt.X
		}
		if pt.Y > max.Y {
			max.Y = pt.Y
		}
	}
	max.X++
	max.Y++

	fmt.Println()
	for y := min.Y; y < max.Y; y++ {
		for x := min.X; x < max.X; x++ {
			pt := image.Pt(x, y)

			chr := " "
			if pt.Eq(image.ZP) {
				chr = "*"
			} else if pt.Eq(droidPos) {
				switch droidDir {
				case north:
					chr = "A"
				case south:
					chr = "V"
				case west:
					chr = "<"
				case east:
					chr = ">"
				}
			} else {
				switch grid[pt] {
				case travel:
					chr = "."
				case wall:
					chr = "#"
				case oxygen:
					chr = "O"
				case path:
					chr = "+"
				}
			}

			fmt.Print(chr)
		}
		fmt.Println()
	}
}

func expandOxygen(grid map[image.Point]content) int {
	var rounds int
	for {
		// until all travel has been replaced with oxygen
		travels := count(grid, travel)
		if travels == 0 {
			return rounds
		}

		rounds++

		oxygens := find(grid, oxygen)
		for _, o := range oxygens {
			for _, cpt := range []image.Point{
				o.Add(image.Pt(0, -1)), // north
				o.Add(image.Pt(0, 1)),  // south
				o.Add(image.Pt(-1, 0)), // west
				o.Add(image.Pt(1, 0)),  // east
			} {
				if grid[cpt] == travel {
					grid[cpt] = oxygen
				}
			}
		}

	}
}

func find(grid map[image.Point]content, n content) []image.Point {
	var out []image.Point
	for pt, con := range grid {
		if con == n {
			out = append(out, pt)
		}
	}
	return out
}

func count(grid map[image.Point]content, n content) int {
	return len(find(grid, n))
}
