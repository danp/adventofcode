package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	st := parse(strings.TrimSpace(string(b)))
	for i := 0; i < 200; i++ {
		st.step()
	}
	fmt.Println("final", st.bugCount())
}

type content int

const (
	empty content = 0
	bug   content = 1
)

type state struct {
	g map[levelPoint]content
}

func (s *state) step() {
	newg := make(map[levelPoint]content)
	var newpts []levelPoint

	for pt, con := range s.g {
		if pt.pt.Eq(cpt) {
			continue
		}

		var bc int
		for apt, acon := range s.adj(pt) {
			if acon == bug {
				bc++
			}
			if _, ok := s.g[apt]; !ok {
				newpts = append(newpts, apt)
			}
		}

		ncon := empty
		if con == bug && bc == 1 {
			ncon = bug
		} else if con == empty && (bc == 1 || bc == 2) {
			ncon = bug
		}

		newg[pt] = ncon
	}

	for _, pt := range newpts {
		var bc int
		for _, acon := range s.adj(pt) {
			if acon == bug {
				bc++
			}
		}

		ncon := empty
		if bc == 1 || bc == 2 {
			ncon = bug
		}

		newg[pt] = ncon
	}

	s.g = newg
}

func (s *state) bugCount() int {
	var bc int
	for _, con := range s.g {
		if con == bug {
			bc++
		}
	}
	return bc
}

type levelPoint struct {
	pt    image.Point
	level int
}

func lp(pt image.Point, level int) levelPoint {
	return levelPoint{pt: pt, level: level}
}

func (l levelPoint) Sub(pt image.Point) levelPoint {
	return lp(l.pt.Sub(pt), l.level)
}

func (l levelPoint) Add(pt image.Point) levelPoint {
	return lp(l.pt.Add(pt), l.level)
}

var cpt = image.Pt(2, 2)

func (s *state) adj(pt levelPoint) map[levelPoint]content {
	out := make(map[levelPoint]content)

	var apts []levelPoint

	above := pt.Sub(image.Pt(0, 1))
	if above.pt.Y == -1 {
		apts = append(apts, lp(image.Pt(2, 1), pt.level-1))
	} else if above.pt.Eq(cpt) {
		for x := 0; x < 5; x++ {
			apts = append(apts, lp(image.Pt(x, 4), pt.level+1))
		}
	} else {
		apts = append(apts, above)
	}

	below := pt.Add(image.Pt(0, 1))
	if below.pt.Y == 5 {
		apts = append(apts, lp(image.Pt(2, 3), pt.level-1))
	} else if below.pt.Eq(cpt) {
		for x := 0; x < 5; x++ {
			apts = append(apts, lp(image.Pt(x, 0), pt.level+1))
		}
	} else {
		apts = append(apts, below)
	}

	right := pt.Add(image.Pt(1, 0))
	if right.pt.X == 5 {
		apts = append(apts, lp(image.Pt(3, 2), pt.level-1))
	} else if right.pt.Eq(cpt) {
		for y := 0; y < 5; y++ {
			apts = append(apts, lp(image.Pt(0, y), pt.level+1))
		}
	} else {
		apts = append(apts, right)
	}

	left := pt.Sub(image.Pt(1, 0))
	if left.pt.X == -1 {
		apts = append(apts, lp(image.Pt(1, 2), pt.level-1))
	} else if left.pt.Eq(cpt) {
		for y := 0; y < 5; y++ {
			apts = append(apts, lp(image.Pt(4, y), pt.level+1))
		}
	} else {
		apts = append(apts, left)
	}

	for _, apt := range apts {
		out[apt] = s.g[apt]
	}

	return out
}

func (s *state) String() string {
	var os strings.Builder

	var minl, maxl int
	for pt := range s.g {
		if pt.level < minl {
			minl = pt.level
		}
		if pt.level > maxl {
			maxl = pt.level
		}
	}

	for l := minl; l <= maxl; l++ {
		fmt.Fprintln(&os, "Level", l)
		for y := 0; y < 5; y++ {
			for x := 0; x < 5; x++ {
				pt := lp(image.Pt(x, y), l)
				switch s.g[pt] {
				case empty:
					os.WriteString(".")
				case bug:
					os.WriteString("#")
				}
			}
			fmt.Fprintln(&os)
		}
		fmt.Fprintln(&os)
	}

	return os.String()
}

func parse(input string) *state {
	grid := make(map[levelPoint]content)
	var pos levelPoint

	for _, l := range strings.Split(input, "\n") {
		for _, c := range l {
			switch c {
			case '.':
				grid[pos] = empty
			case '#':
				grid[pos] = bug
			}
			pos.pt.X++
		}
		pos.pt.X = 0
		pos.pt.Y++
	}

	return &state{g: grid}
}
