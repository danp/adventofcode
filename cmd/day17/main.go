package main

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

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

	grid, _, rect := findGrid(program)

	if os.Args[1] == "part1" {
		display(grid, rect)
		fmt.Println()

		inters := make(map[image.Point]bool)
		for pt, con := range grid {
			if con != scaffold {
				continue
			}

			if isIntersection(grid, pt) {
				inters[pt] = true
			}
		}

		for pt := range inters {
			grid[pt] = intersection
		}
		display(grid, rect)

		var isum int
		for pt := range inters {
			p := pt.X * pt.Y
			isum += p
		}
		fmt.Println(isum)
	}

	display(grid, rect)

	// pf := pathFinder{grid: grid, rs: rs}
	// pf.run()
	// fmt.Println(pf.main)

	pairs := map[string]string{
		"V": "R,12",
		"W": "L,8",
		"X": "L,4",
		"Y": "R,6",
		"Z": "L,6",
	}

	functions := map[string]string{
		"A": "VWXX",
		"B": "WYZ",
		"C": "WXVZX",
	}

	mrout := "ABABCACACB"

	var buf bytes.Buffer

	// main routine
	for i, c := range mrout {
		if i > 0 {
			fmt.Fprint(&buf, ",")
		}
		fmt.Fprint(&buf, string(c))
	}
	fmt.Fprintln(&buf)

	// functions
	for _, fn := range []string{"A", "B", "C"} {
		for i, pc := range functions[fn] {
			if i > 0 {
				fmt.Fprint(&buf, ",")
			}
			fmt.Fprint(&buf, pairs[string(pc)])
		}
		fmt.Fprintln(&buf)
	}

	// video feed
	fmt.Fprintln(&buf, "n")

	fmt.Printf("%q\n", buf.String())

	input := func() (int, error) {
		c, err := buf.ReadByte()
		return int(c), err
	}

	output := func(x int) error {
		if x > 255 {
			fmt.Println(x)
		}
		return nil
	}

	program[0] = 2
	mem := make([]int, 4096)

	if err := intcode.Run(program, mem, input, output); err != nil {
		panic(err)
	}
}

type command struct {
	t string
	n int
}

func cs(commands []command) string {
	var m int
	var out []string
	for _, c := range commands {
		if c.t != "" {
			if m > 0 {
				out = append(out, strconv.Itoa(m))
				m = 0
			}
			out = append(out, c.t)
		} else {
			m += c.n
		}
	}
	if m > 0 {
		out = append(out, strconv.Itoa(m))
	}
	return strings.Join(out, ",")
}

type routine struct {
	n string
	c []command
}

type pathFinder struct {
	grid map[image.Point]content
	rs   rstate

	main []routine
}

func (p pathFinder) run() {
	var nscaf int
	var rpt image.Point
	for pt, con := range p.grid {
		if con == robot {
			nscaf++ // robot is always on scaffold
			rpt = pt
		}
		if con == scaffold {
			nscaf++
		}
	}

	rs := p.rs

	visited := map[image.Point]bool{
		rpt: true,
	}

	var seq []command

	for {
		sp := findPath(
			p.grid,
			rpt,
			rs,
			func(p []image.Point) bool {
				return !visited[p[len(p)-1]]
			},
			func(pt image.Point) bool {
				return p.grid[pt] == scaffold
			},
		)

		if len(sp) < 2 {
			if len(visited) != nscaf {
				panic("didn't visit everything")
			}
			break
		}

		np := sp[1]

		var ns rstate
		if np.X < rpt.X {
			ns = left
		} else if np.X > rpt.X {
			ns = right
		} else if np.Y < rpt.Y {
			ns = up
		} else if np.Y > rpt.Y {
			ns = down
		}

		if rs != ns {
			if d := ns - rs; d == (right-up) || d == (down-right) || d == (left-down) || d == (up-left) || d == (left-right) || d == (up-down) || d == (down-up) {
				seq = append(seq, command{t: "R"})
				rs++
			} else if d == (left-up) || d == (down-left) || d == (right-down) || d == (up-right) {
				seq = append(seq, command{t: "L"})
				rs--
			} else {
				panic("what")
			}

			if rs > left {
				rs = up
			} else if rs < up {
				rs = left
			}
		} else {
			seq = append(seq, command{n: 1})
			rpt = np
		}

		visited[rpt] = true
	}

	fmt.Println(cs(seq))

	cparts := strings.Split(cs(seq), ",")
	pairs := make(map[string]byte)
	pn := byte('V')

	for i := 0; i < len(cparts); i += 2 {
		cp := strings.Join(cparts[i:i+2], ",")
		if _, ok := pairs[cp]; !ok {
			pairs[cp] = pn
			pn++
		}

		fmt.Print(string(pairs[cp]))
	}
	fmt.Println()
	fmt.Println()

	for s, pn := range pairs {
		fmt.Println(s, string(pn))
	}
}

func findGrid(program []int) (map[image.Point]content, rstate, image.Rectangle) {
	grid := make(map[image.Point]content)
	pt := image.ZP
	var maxPt image.Point
	var rs rstate

	output := func(x int) error {
		switch byte(x) {
		case '#':
			grid[pt] = scaffold
		case '.':
			grid[pt] = open
		case '\n':
			maxPt.Y = pt.Y
			pt.Y++
			pt.X = 0
			return nil
		case '^':
			grid[pt] = robot
			rs = up
		case 'v':
			grid[pt] = robot
			rs = down
		case '>':
			grid[pt] = robot
			rs = right
		case '<':
			grid[pt] = robot
			rs = left
		}

		pt.X++
		if pt.X > maxPt.X {
			maxPt.X = pt.X
		}

		return nil
	}

	mem := make([]int, 4096)
	if err := intcode.Run(program, mem, nil, output); err != nil {
		panic(err)
	}

	return grid, rs, image.Rect(0, 0, maxPt.X, maxPt.Y)
}

func display(grid map[image.Point]content, rect image.Rectangle) {
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			pt := image.Pt(x, y)
			con := grid[pt]

			switch con {
			case scaffold:
				fmt.Print("#")
			case open:
				fmt.Print(".")
			case robot:
				fmt.Print("R")
			case intersection:
				fmt.Print("O")
			}
		}
		fmt.Println()
	}
}

func isIntersection(grid map[image.Point]content, pt image.Point) bool {
	// 0 1 2
	// 3 x 4
	// 5 6 7

	var adjacent []content
	for y := pt.Y - 1; y <= pt.Y+1; y++ {
		for x := pt.X - 1; x <= pt.X+1; x++ {
			apt := image.Pt(x, y)
			if apt.Eq(pt) {
				continue
			}

			adjacent = append(adjacent, grid[apt])
		}
	}

	return adjacent[1] == scaffold && adjacent[3] == scaffold && adjacent[4] == scaffold && adjacent[6] == scaffold
}

type content int

const (
	unknown  content = 0
	scaffold content = 1
	open     content = 2
	robot    content = 3

	intersection content = 99
)

type rstate int

const (
	up    rstate = 1
	right rstate = 2
	down  rstate = 3
	left  rstate = 4

	spaced = 99
)

// findPath finds the shortest path from start to where check returns true.
func findPath(grid map[image.Point]content, start image.Point, rs rstate, check func(path []image.Point) bool, consider func(pt image.Point) bool) []image.Point {
	seen := make(map[image.Point]bool)

	q := [][]image.Point{{start}}
	for len(q) > 0 {
		path := q[0]
		q = q[1:]

		if check(path) {
			return path
		}

		pt := path[len(path)-1]

		poss := []image.Point{
			pt.Add(image.Pt(0, -1)), // north
			pt.Add(image.Pt(1, 0)),  // east
			pt.Add(image.Pt(0, 1)),  // south
			pt.Add(image.Pt(-1, 0)), // west
		}

		for possidx := rs; possidx < rs+left; possidx++ {
			pi := possidx
			if pi > left {
				pi -= left
			}
			cpt := poss[pi-1]

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
