package main

import (
	"fmt"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	var scanners []*scanner
	var s *scanner
	for _, l := range lines {
		if l == "" {
			if s != nil {
				scanners = append(scanners, s)
			}
			continue
		}
		if strings.HasPrefix(l, "---") {
			s = newScanner(strings.Fields(l)[2])
			continue
		}

		ns := strings.Split(l, ",")
		pt := point3{scaffold.Int(ns[0]), scaffold.Int(ns[1]), scaffold.Int(ns[2])}
		s.beacons[pt] = struct{}{}

	}
	if s != nil {
		scanners = append(scanners, s)
	}

	pt1 := point3{-618, -824, -621}
	fmt.Printf("pt1: %v\n", pt1)
	// off := point3{68, -1246, -43}

	pt2 := point3{686, 422, 578}
	fmt.Printf("pt2: %v\n", pt2)

	for _, o := range pt2.orientations() {
		fmt.Printf("o: %v sub: %v\n", o, pt1.sub(o))
	}

	// given 1,2,3
	// if roll forward
	// becomes 1,3,2

	// given -1,-2,-3
	// if roll forward
	// becomes -1,-3,-2

	// what translations are needed to make this work?

}

type point3 struct {
	x, y, z int
}

func (p point3) sub(o point3) point3 {
	return point3{p.x - o.x, p.y - o.y, p.z - o.z}
}

func (p point3) roll() point3 {
	return point3{p.x, p.z, -p.y}
}

func (p point3) turn() point3 {
	return point3{-p.y, p.x, p.z}
}

func (p point3) orientations() []point3 {
	var out []point3
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			p = p.roll()
			out = append(out, p)
			for l := 0; l < 3; l++ {
				p = p.turn()
				out = append(out, p)
			}
		}
		p = p.roll().turn().roll()
	}
	return out
}

type scanner struct {
	id      string
	beacons map[point3]struct{}
}

func newScanner(id string) *scanner {
	return &scanner{id: id, beacons: make(map[point3]struct{})}
}
