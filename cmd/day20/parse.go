package main

import (
	"fmt"
	"image"
	"regexp"
	"sort"
	"strings"
)

func parse(input string) map[image.Point]content {
	lines := strings.Split(input, "\n")
	// 0: collect top letters, 1: in maze, 2: collect bottom letters
	var (
		st   int
		xoff int
		maxX int
	)

	grid := make(map[image.Point]content)
	var pos image.Point

	vletters := make(map[int]string)
	for _, l := range lines {
		switch st {
		case 0:
			if fw := strings.Index(l, "#"); fw >= 0 {
				xoff = fw
				maxX = strings.LastIndex(l, "#") - fw

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
				mpt := image.Pt((idx - xoff), pos.Y)
				pt := ppt(mpt)
				grid[pt] = content{portal: ls, pinner: mpt.X > 0 && mpt.X < maxX}
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

				// vletters here are always inner
				grid[ppt(pt)] = content{portal: ls, pinner: true}
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
				desc := fmt.Sprintf("%s %s", pcon.portal, ppt)
				if pcon.pinner {
					desc += "+"
				}
				pmap[pnum] = desc
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
