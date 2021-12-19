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

	type mapping struct {
		idx int
		tf  func(point3) point3
		off point3
	}
	mappings := make(map[int]mapping)

	for s1i, s1 := range scanners {
		for s2i, s2 := range scanners {
			if s1i == s2i {
				continue
			}

			for tfi, tf := range translations() {
				offs := make(map[point3]int)
				for s1pt := range s1.beacons {
					for s2pt := range s2.beacons {
						offs[s1pt.sub(tf(s2pt))]++
					}
				}
				for off, v := range offs {
					if v >= 12 {
						fmt.Printf("tfi: %v off: %v v: %v\n", tfi, off, v)
						mappings[s1i] = mapping{s2i, tf, off}
					}
				}
			}
		}
	}

	fmt.Printf("mappings: %v\n", mappings)
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

func translations() []func(point3) point3 {
	var out []func(point3) point3
	var chain []func(point3) point3
	emit := func() {
		cc := make([]func(point3) point3, len(chain))
		copy(cc, chain)
		out = append(out, func(p point3) point3 {
			for _, c := range cc {
				p = c(p)
			}
			return p
		})
	}
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			chain = append(chain, (point3).roll)
			emit()
			for l := 0; l < 3; l++ {
				chain = append(chain, (point3).turn)
				emit()
			}
		}
		chain = append(chain, (point3).roll, (point3).turn, (point3).roll)
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
