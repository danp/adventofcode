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

	mappings := make(map[int]map[int]mapping)

	for s1i, s1 := range scanners {
		for s2i, s2 := range scanners {
			if s1i == s2i {
				continue
			}

			for _, tf := range translations() {
				offs := make(map[point3]int)
				for s1pt := range s1.beacons {
					for s2pt := range s2.beacons {
						offs[s1pt.sub(tf(s2pt))]++
					}
				}
				for off, v := range offs {
					if v >= 12 {
						if mappings[s1i] == nil {
							mappings[s1i] = make(map[int]mapping)
						}
						mappings[s1i][s2i] = mapping{s2i, tf, off}
					}
				}
			}
		}
	}

	beacons := make(map[point3]struct{})
	for si, s := range scanners {
		mf := mt(mappings, si)
		for pt := range s.beacons {
			npt := mf(pt)
			beacons[npt] = struct{}{}
		}
	}

	fmt.Printf("len(beacons): %v\n", len(beacons))
}

type mapping struct {
	idx int
	tf  func(point3) point3
	off point3
}

func mt(mappings map[int]map[int]mapping, src int) func(point3) point3 {
	if src == 0 {
		return func(p point3) point3 {
			return p
		}
	}

	q := [][]int{{0}}
	for len(q) > 0 {
		p := q[0]
		q = q[1:]

		last := p[len(p)-1]
		if last == src {
			var chain []func(point3) point3
			for i := len(p) - 1; i > 0; i-- {
				m := mappings[p[i-1]][p[i]]
				chain = append(chain, func(p point3) point3 {
					return m.tf(p).add(m.off)
				})
			}
			return func(p point3) point3 {
				for _, f := range chain {
					p = f(p)
				}
				return p
			}
		}

		for n := range mappings[last] {
			newp := make([]int, len(p))
			copy(newp, p)
			newp = append(newp, n)
			q = append(q, newp)
		}
	}

	panic("no")
}

type point3 struct {
	x, y, z int
}

func (p point3) sub(o point3) point3 {
	return point3{p.x - o.x, p.y - o.y, p.z - o.z}
}

func (p point3) add(o point3) point3 {
	return point3{p.x + o.x, p.y + o.y, p.z + o.z}
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
