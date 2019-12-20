package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	grid := parse(string(b))
	fmt.Println(solve(grid))
}

type content struct {
	portal string
	wall   bool
	travel bool
}

func parse(input string) map[image.Point]content {
	lines := strings.Split(input, "\n")
	// 0: collect top letters, 1: in maze, 2: collect bottom letters
	var (
		st   int
		xoff int
	)

	grid := make(map[image.Point]content)
	var pos image.Point

	vletters := make(map[int]string)
	for _, l := range lines {
		switch st {
		case 0:
			if fw := strings.Index(l, "#"); fw >= 0 {
				xoff = fw

				for idx, ls := range vletters {
					pt := ppt(image.Pt(idx-xoff, pos.Y))
					grid[pt] = content{portal: ls}
				}

				vletters = make(map[int]string)

				st++
			}
		case 1:
			if strings.Index(l, "#") == -1 {
				st++
			}
		}

		switch st {
		case 0:
			collectVLetters(vletters, l)
		case 1:
			parseLine(grid, pos, l)

			for idx, ls := range collectHLetters(l) {
				pt := ppt(image.Pt((idx - xoff), pos.Y))
				grid[pt] = content{portal: ls}
			}

			collectVLetters(vletters, l)
			var del []int
			for idx, ls := range vletters {
				if len(ls) < 2 {
					continue
				}

				pt := image.Pt((idx - xoff), pos.Y)

				// find where it came from, we are either +2 from
				// a travel or we haven't gotten there yet
				cpt := pt.Sub(image.Pt(0, 2))
				if con, ok := grid[cpt]; ok && con.travel {
					pt = cpt
				} else {
					pt.Y += 1
				}

				grid[ppt(pt)] = content{portal: ls}
				del = append(del, idx)
			}

			for _, d := range del {
				delete(vletters, d)
			}

			pos.Y++
		case 2:
			collectVLetters(vletters, l)
		}
	}

	for idx, ls := range vletters {
		pt := ppt(image.Pt((idx - xoff), pos.Y-1))
		grid[pt] = content{portal: ls}
	}

	return grid
}

func display(grid map[image.Point]content) string {
	var out string
	var maxX int
	var pnum int
	pmap := make(map[int]string)

	for y := 0; ; y++ {
		if _, ok := grid[image.Pt(0, y)]; !ok {
			break
		}

		var lout string

		for x := 0; ; x++ {
			pt := image.Pt(x, y)
			con, ok := grid[pt]
			if !ok && x > maxX {
				break
			}
			if pt.X > maxX {
				maxX = pt.X
			}

			ppt := ppt(pt)
			if pcon, ok := grid[ppt]; ok {
				pmap[pnum] = fmt.Sprintf("%s %s", pcon.portal, ppt)
				lout += pnumchr(pnum)
				pnum++
			} else if con.wall {
				lout += "#"
			} else if con.travel {
				lout += "."
			} else {
				lout += " "
			}
		}

		out += lout + "\n"
	}

	pmk := make([]int, 0, len(pmap))
	for pnum := range pmap {
		pmk = append(pmk, pnum)
	}
	sort.Ints(pmk)
	for _, pnum := range pmk {
		out += fmt.Sprintln(pnumchr(pnum), pmap[pnum])
	}

	return out
}

func pnumchr(pnum int) string {
	chr := 'A' + pnum
	if chr > 'Z' {
		chr = 'a' + chr - 'Z' - 1
	}
	return string(chr)
}

// translate maze point pt into its associated portal point
func ppt(pt image.Point) image.Point {
	return pt.Add(image.Pt(1, 1)).Mul(-1)
}

// translate portal point pt it its associated maze point
func mpt(pt image.Point) image.Point {
	return pt.Mul(-1).Sub(image.Pt(1, 1))
}

var vlre = regexp.MustCompile(`[A-Z]`)

func collectVLetters(letters map[int]string, l string) {
	for _, m := range vlre.FindAllStringIndex(l, -1) {
		if m[0] == 0 {
			continue
		}
		if l[m[0]-1] == ' ' && (m[1] == len(l) || l[m[1]] == ' ') {
			lts := letters[m[0]]
			lts += l[m[0]:m[1]]
			letters[m[0]] = lts
		}
	}
}

func collectHLetters(l string) map[int]string {
	out := make(map[int]string)

	ms := strings.IndexAny(l, "#.")
	if bm := strings.TrimSpace(l[:ms]); bm != "" {
		out[ms] = bm
	}

	me := strings.LastIndexAny(l, "#.")
	if am := strings.TrimSpace(l[me+1:]); am != "" {
		out[me] = am
	}

	hs := strings.IndexFunc(l[ms:], func(c rune) bool { return c != '.' && c != '#' })
	if hs == -1 || hs+ms == me {
		return out
	}
	hs += ms
	he := strings.IndexAny(l[hs:], "#.") + hs - 1

	if l[hs-1] == '.' {
		out[hs-1] = l[hs : hs+2]
	}

	if l[he+1] == '.' {
		out[he+1] = l[he-1 : he+1]
	}

	return out
}

func parseLine(grid map[image.Point]content, pos image.Point, l string) {
	start := strings.IndexAny(l, "#.")

	for _, c := range l[start:] {
		switch c {
		case '#':
			grid[pos] = content{wall: true}
		case '.':
			grid[pos] = content{travel: true}
		}
		pos = pos.Add(image.Pt(1, 0))
	}
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
			if grid[pt].travel {
				return true
			}
			if grid[pt].portal != "" {
				return true
			}
			return false
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
