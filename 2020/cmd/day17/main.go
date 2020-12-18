package main

import (
	"fmt"
	"image"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	do3(lines)
	do4(lines)
}

func do3(lines []string) {
	world := make(map[point3]bool)

	for y, l := range lines {
		for x, c := range l {
			var active bool
			if c == '#' {
				active = true
			}
			pt := point3{Pt: image.Pt(x, y)}
			world[pt] = active
		}
	}

	for i := 0; i < 6; i++ {
		nw := copyWorld3(world)
		disc := make(map[point3]bool)

		for p3, act := range world {
			var nact int
			for _, n := range p3.neighbors() {
				nearActive, ok := world[n]
				if nearActive {
					nact++
				}
				if !ok {
					disc[n] = true
				}
			}

			nwact := act
			if act && (nact < 2 || nact > 3) {
				nwact = false
			} else if !act && nact == 3 {
				nwact = true
			}
			nw[p3] = nwact
		}

		for p3 := range disc {
			var nact int
			for _, n := range p3.neighbors() {
				nearActive, _ := world[n]
				if nearActive {
					nact++
				}
			}

			nwact := false
			if nact == 3 {
				nwact = true
			}
			nw[p3] = nwact
		}

		world = nw
	}

	var numAct int
	for _, act := range world {
		if act {
			numAct++
		}
	}
	fmt.Println(numAct)
}

func do4(lines []string) {
	world := make(map[point4]bool)

	for y, l := range lines {
		for x, c := range l {
			var active bool
			if c == '#' {
				active = true
			}
			pt := point4{Pt: image.Pt(x, y)}
			world[pt] = active
		}
	}

	for i := 0; i < 6; i++ {
		nw := copyWorld4(world)
		disc := make(map[point4]bool)

		for pt, act := range world {
			var nact int
			for _, n := range pt.neighbors() {
				nearActive, ok := world[n]
				if nearActive {
					nact++
				}
				if !ok {
					disc[n] = true
				}
			}

			nwact := act
			if act && (nact < 2 || nact > 3) {
				nwact = false
			} else if !act && nact == 3 {
				nwact = true
			}
			nw[pt] = nwact
		}

		for pt := range disc {
			var nact int
			for _, n := range pt.neighbors() {
				nearActive, _ := world[n]
				if nearActive {
					nact++
				}
			}

			nwact := false
			if nact == 3 {
				nwact = true
			}
			nw[pt] = nwact
		}

		world = nw
	}

	var numAct int
	for _, act := range world {
		if act {
			numAct++
		}
	}
	fmt.Println(numAct)
}

type point3 struct {
	Pt image.Point
	Z  int
}

func (p point3) neighbors() []point3 {
	var out []point3
	for z := p.Z - 1; z <= p.Z+1; z++ {
		for y := p.Pt.Y - 1; y <= p.Pt.Y+1; y++ {
			for x := p.Pt.X - 1; x <= p.Pt.X+1; x++ {
				if z == p.Z && y == p.Pt.Y && x == p.Pt.X {
					continue
				}
				out = append(out, point3{Pt: image.Pt(x, y), Z: z})
			}
		}
	}
	return out
}

func copyWorld3(w map[point3]bool) map[point3]bool {
	out := make(map[point3]bool)
	for k, v := range w {
		out[k] = v
	}
	return out
}

type point4 struct {
	Pt image.Point
	Z  int
	W  int
}

func (p point4) neighbors() []point4 {
	var out []point4
	for w := p.W - 1; w <= p.W+1; w++ {
		for z := p.Z - 1; z <= p.Z+1; z++ {
			for y := p.Pt.Y - 1; y <= p.Pt.Y+1; y++ {
				for x := p.Pt.X - 1; x <= p.Pt.X+1; x++ {
					if w == p.W && z == p.Z && y == p.Pt.Y && x == p.Pt.X {
						continue
					}
					out = append(out, point4{Pt: image.Pt(x, y), Z: z, W: w})
				}
			}
		}
	}
	return out
}

func copyWorld4(w map[point4]bool) map[point4]bool {
	out := make(map[point4]bool)
	for k, v := range w {
		out[k] = v
	}
	return out
}
