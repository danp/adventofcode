package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	grid := parse(string(b))
	if os.Args[1] == "part1" {
		fmt.Println(solve(grid))
	} else {
		fmt.Println(solveR(grid))
	}
}

type content struct {
	portal string
	pinner bool
	wall   bool
	travel bool
}

// translate maze point pt into its associated portal point
func ppt(pt image.Point) image.Point {
	return pt.Add(image.Pt(1, 1)).Mul(-1)
}

// translate portal point pt it its associated maze point
func mpt(pt image.Point) image.Point {
	return pt.Mul(-1).Sub(image.Pt(1, 1))
}

func bpt(pt image.Point) image.Point {
	n := level(pt) * 1000

	if pt.X < 0 {
		return pt.Add(image.Pt(n, n))
	} else {
		return pt.Sub(image.Pt(n, n))
	}
}

func level(pt image.Point) int {
	if pt.X < 0 {
		pt.X *= -1
	}
	return pt.X / 1000
}

func lpt(pt image.Point, level int) image.Point {
	n := level * 1000
	if pt.X < 0 {
		return pt.Sub(image.Pt(n, n))
	} else {
		return pt.Add(image.Pt(n, n))
	}
}

// findPathR finds the shortest path from start to where check returns true, recursively.
func findPathR(grid map[image.Point]content, start image.Point, check func(path []image.Point) bool, consider func(pt image.Point) bool) []image.Point {
	pmap := makePortalMap(grid)
	seen := make(map[image.Point]bool)

	q := [][]image.Point{{start}}
	for len(q) > 0 {
		path := q[0]
		q = q[1:]

		if check(path) {
			return path
		}

		pt := path[len(path)-1]
		level := level(pt)

		moves := []image.Point{
			pt.Add(image.Pt(0, -1)), // north
			pt.Add(image.Pt(0, 1)),  // south
			pt.Add(image.Pt(-1, 0)), // west
			pt.Add(image.Pt(1, 0)),  // east
		}

		poss := make([]image.Point, 0, len(moves))
		// can we step to these places given the level?
		for _, cpt := range moves {
			con, ok := grid[ppt(bpt(cpt))]
			if !ok {
				poss = append(poss, cpt)
				continue
			}

			if level == 0 {
				if con.portal != "" {
					if !con.pinner && con.portal != "AA" && con.portal != "ZZ" {
						continue
					}
				}
			} else if level > 0 {
				if con.portal != "" && (con.portal == "AA" || con.portal == "ZZ") {
					continue
				}
			}
			poss = append(poss, cpt)
		}

		// if we are at a portal, add its other side with necessary level change
		ppt := ppt(bpt(pt))
		os, ok := pmap[ppt]
		if ok && !bpt(path[len(path)-2]).Eq(os) { // don't immediately go back through the portal
			pcon := grid[ppt]
			if pcon.pinner {
				os = lpt(os, level+1)
				if level < 30 {
					poss = append(poss, os)
				} else {
					fmt.Println(pt, "already at level", level, "not considering deeper")
				}
			} else {
				os = lpt(os, level-1)
				poss = append(poss, os)
			}
		}

		for _, cpt := range poss {
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

// findPath finds the shortest path from start to where check returns true.
func findPath(grid map[image.Point]content, start image.Point, check func(path []image.Point) bool, consider func(pt image.Point) bool) []image.Point {
	pmap := makePortalMap(grid)
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
			pt.Add(image.Pt(0, 1)),  // south
			pt.Add(image.Pt(-1, 0)), // west
			pt.Add(image.Pt(1, 0)),  // east
		}

		ppt := ppt(pt)
		if os, ok := pmap[ppt]; ok {
			poss = append(poss, os)
		}

		for _, cpt := range poss {
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

func solveR(grid map[image.Point]content) int {
	var aa, zz image.Point
	for pt, con := range grid {
		if con.portal == "AA" {
			aa = mpt(pt)
		}
		if con.portal == "ZZ" {
			zz = mpt(pt)
		}
	}
	if aa.Eq(image.ZP) || zz.Eq(image.ZP) {
		panic("aa or zz zero")
	}

	sp := findPathR(
		grid,
		aa,
		func(p []image.Point) bool {
			return p[len(p)-1].Eq(zz)
		},
		func(pt image.Point) bool {
			return grid[bpt(pt)].travel
		},
	)

	return len(sp) - 1
}

func solve(grid map[image.Point]content) int {
	var aa, zz image.Point
	for pt, con := range grid {
		if con.portal == "AA" {
			aa = mpt(pt)
		}
		if con.portal == "ZZ" {
			zz = mpt(pt)
		}
	}
	if aa.Eq(image.ZP) || zz.Eq(image.ZP) {
		panic("aa or zz zero")
	}

	sp := findPath(
		grid,
		aa,
		func(p []image.Point) bool {
			return p[len(p)-1].Eq(zz)
		},
		func(pt image.Point) bool {
			return grid[pt].travel
		},
	)

	return len(sp) - 1
}

// portal at portal point key goes to travel at maze point value
func makePortalMap(grid map[image.Point]content) map[image.Point]image.Point {
	t := make(map[string][]image.Point)
	for pt, con := range grid {
		p := con.portal
		if p == "" || p == "AA" || p == "ZZ" {
			continue
		}
		x := t[p]
		if x == nil {
			x = make([]image.Point, 0, 1)
		}
		x = append(x, pt)
		t[p] = x
	}

	out := make(map[image.Point]image.Point)
	for _, pts := range t {
		if len(pts) != 2 {
			panic("not 2 pts")
		}
		out[pts[0]] = mpt(pts[1])
		out[pts[1]] = mpt(pts[0])
	}
	return out
}
