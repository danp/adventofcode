package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	linestrings := strings.Split(strings.TrimSpace(string(b)), "\n")

	var lines []line
	for _, ls := range linestrings {
		l, err := parseLine(ls)
		if err != nil {
			panic(err)
		}
		lines = append(lines, l)
	}

	if len(lines)%2 != 0 {
		panic("non-even number of lines")
	}

	switch os.Args[1] {
	case "points":
		for i := 0; i < len(lines); i += 2 {
			fmt.Println(minPointDistance(lines[i].intersections(lines[i+1])))
		}
	case "steps":
		for i := 0; i < len(lines); i += 2 {
			fmt.Println(minPointSteps(lines[i].pointSteps(), lines[i+1].pointSteps()))
		}
	}
}

type dir int

func (d dir) String() string {
	switch d {
	case up:
		return "up"
	case right:
		return "right"
	case down:
		return "down"
	case left:
		return "left"
	default:
		return "unknown"
	}
}

const (
	up    dir = 1
	right dir = 2
	down  dir = 3
	left  dir = 4
)

type move struct {
	d dir
	c int
}

type line struct {
	moves []move
}

type point struct{ x, y int }

func (l line) points() map[point]bool {
	out := make(map[point]bool)
	var p point

	for _, m := range l.moves {
		for i := 0; i < m.c; i++ {
			switch m.d {
			case up:
				p.y++
			case down:
				p.y--
			case right:
				p.x++
			case left:
				p.x--
			}

			out[p] = true
		}
	}

	return out
}

func (l line) pointSteps() map[point]int {
	out := make(map[point]int)
	var p point
	var s int

	for _, m := range l.moves {
		for i := 0; i < m.c; i++ {
			s++

			switch m.d {
			case up:
				p.y++
			case down:
				p.y--
			case right:
				p.x++
			case left:
				p.x--
			}

			if _, ok := out[p]; !ok {
				out[p] = s
			}
		}
	}

	return out
}

func (l line) intersections(o line) []point {
	ours := l.points()
	theirs := o.points()

	needles := ours
	haystack := theirs

	if len(theirs) < len(needles) {
		needles = theirs
		haystack = ours
	}

	var out []point
	for p := range needles {
		_, ok := haystack[p]
		if ok {
			out = append(out, p)
		}
	}
	return out
}

func parseLine(input string) (line, error) {
	var l line

	parts := strings.Split(input, ",")
	for _, p := range parts {
		c, err := strconv.Atoi(p[1:len(p)])
		if err != nil {
			return l, err
		}

		var d dir
		switch p[0] {
		case 'U':
			d = up
		case 'R':
			d = right
		case 'D':
			d = down
		case 'L':
			d = left
		default:
			return l, fmt.Errorf("unknown direction in instruction %q", p)
		}

		l.moves = append(l.moves, move{d: d, c: c})
	}

	return l, nil
}

func minPointDistance(points []point) int {
	dist := -1

	for _, p := range points {
		x := p.x
		y := p.y
		if x < 0 {
			x *= -1
		}
		if y < 0 {
			y *= -1
		}

		d := x + y
		if dist == -1 || d < dist {
			dist = d
		}
	}

	return dist
}

func minPointSteps(p1 map[point]int, p2 map[point]int) int {
	steps := -1

	needles := p1
	haystack := p2
	if len(haystack) < len(needles) {
		needles, haystack = haystack, needles
	}

	for n, ns := range needles {
		hs, ok := haystack[n]
		if !ok {
			continue
		}

		ts := ns + hs
		if steps == -1 || ts < steps {
			steps = ts
		}
	}

	return steps
}
