package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

var linelen int

func main() {
	lines := scaffold.Lines()

	var tiles []tile

	var nt tile
	for _, l := range lines {
		if l == "" {
			tiles = append(tiles, nt)
			nt = tile{}
			continue
		}

		if strings.HasPrefix(l, "Tile") {
			nt.id = l[5 : len(l)-1]
			continue
		}

		nt.data += l
		if len(l) > linelen {
			linelen = len(l)
		}
	}
	tiles = append(tiles, nt)
	sideLen := int(math.Sqrt(float64(len(tiles))))

	allFrames := make(map[string]frame)
	for _, t0 := range tiles {
		t0 := t0
		for _, t0f := range t0.frames() {
			t0f := t0f
			allFrames[t0f.id()] = t0f
		}
	}

	pairings := map[int]int{
		0: 2, // top and bottom
		2: 0, // bottom and top
		1: 3, // right and left
		3: 1, // left and right
	}

	joinings := make(map[string]map[int]map[string]bool)
	for _, t0 := range tiles {
		for _, t0f := range t0.frames() {
			if joinings[t0f.id()] == nil {
				joinings[t0f.id()] = make(map[int]map[string]bool)
			}

			for _, t1 := range tiles {
				if t0.id == t1.id {
					continue
				}

				for _, t1f := range t1.frames() {
					for pk, pv := range pairings {
						if t1f.edges[pv] == t0f.edges[pk] {
							if joinings[t0f.id()][pk] == nil {
								joinings[t0f.id()][pk] = make(map[string]bool)
							}
							joinings[t0f.id()][pk][t1f.id()] = true
						}
					}
				}
			}
		}
	}

	var pq [][]string
	for _, p := range neighbors(joinings, allFrames, nil, sideLen) {
		pq = append(pq, []string{p})
	}
	var p []string
	for len(pq) > 0 {
		p = pq[0]
		pq = pq[1:]

		if len(p) == len(tiles) {
			break
		}

		for _, n := range neighbors(joinings, allFrames, p, sideLen) {
			newp := make([]string, len(p)+1)
			copy(newp, p)
			newp[len(newp)-1] = n
			pq = append(pq, newp)
		}
	}

	for i := 0; i < len(p); i++ {
		if i > 0 && i%sideLen == 0 {
			fmt.Println()
		}
		fmt.Printf("%s\t", p[i])
	}
	fmt.Println()

	actualImage := make([]string, sideLen*(linelen-2)) // -2 for removing header and footer
	for y := 0; y < sideLen; y++ {
		outs := make([][]string, len(strings.Split(allFrames[p[0]].DataString(), "\n")))
		for i := y * sideLen; i < (y+1)*sideLen; i++ {
			for j, o := range strings.Split(allFrames[p[i]].DataString(), "\n") {
				outs[j] = append(outs[j], o, " ")
			}
		}
		for _, o := range outs {
			fmt.Println(strings.Join(o, ""))
		}

		for i := y * sideLen; i < (y+1)*sideLen; i++ {
			fd := strings.Split(allFrames[p[i]].DataString(), "\n")
			fd = fd[1 : len(fd)-1]
			for ln, l := range fd {
				actualImage[y*len(fd)+ln] += l[1 : len(l)-1]
			}
		}
	}

	cornerProd := 1
	for ci, corner := range []int{0, sideLen - 1, sideLen * (sideLen - 1), sideLen*sideLen - 1} {
		fmt.Println(ci, corner, p[corner], allFrames[p[corner]].t.id)
		cornerProd *= scaffold.Int(allFrames[p[corner]].t.id)
	}
	fmt.Println(cornerProd)

	fmt.Println(strings.Join(actualImage, "\n"))
	fmt.Println()

	tf := []func([]string) []string{
		nil,
		rotates,
		func(s []string) []string { return rotates(rotates(s)) },
		func(s []string) []string { return rotates(rotates(rotates(s))) },
	}

outer:
	for _, otf := range []func([]string) []string{nil, flips} {
		for _, itf := range tf {
			image := actualImage
			if otf != nil {
				image = otf(image)
			}
			if itf != nil {
				image = itf(image)
			}

			fmt.Println(strings.Join(image, "\n"))
			fmt.Println()

			found, smh := findSeaMonsters(image)
			if found > 0 {
				fmt.Println("found", found, "sea monsters consuming", smh, "hashes")

				var hashes int
				for _, l := range image {
					for _, c := range l {
						if c == '#' {
							hashes++
						}
					}
				}
				hashes -= smh
				fmt.Println(hashes, "hashes remaining")
				break outer
			}
		}
	}
}

func neighbors(joinings map[string]map[int]map[string]bool, allFrames map[string]frame, path []string, sideLen int) []string {
	x, y := len(path)%sideLen, len(path)/sideLen
	py := (len(path) - 1) / sideLen

	pathFrames := make([]frame, 0, len(path))
	usedTiles := make(map[string]bool)
	for _, p := range path {
		pf := allFrames[p]
		pathFrames = append(pathFrames, pf)
		usedTiles[pf.t.id] = true
	}

	var cf frame
	if len(pathFrames) > 0 {
		cf = pathFrames[len(pathFrames)-1]
	}

	var out []string

nextJoining:
	for nfid, dirs := range joinings {
		// top row
		if y == 0 && len(dirs[0]) > 0 {
			continue
		}
		// right edge
		if x == sideLen-1 && len(dirs[1]) > 0 {
			continue
		}

		// middle
		if x > 0 && x < sideLen-1 && y > 0 && y < sideLen-1 {
			for _, dv := range dirs {
				if len(dv) == 0 {
					continue nextJoining
				}
			}
		}

		// bottom row
		if y == sideLen-1 && len(dirs[2]) > 0 {
			continue
		}
		// left edge
		if x == 0 && len(dirs[3]) > 0 {
			continue
		}

		nf := allFrames[nfid]
		if usedTiles[nf.t.id] {
			continue
		}

		// short-circuit initial case
		if len(path) == 0 {
			out = append(out, nfid)
			continue
		}

		if y == py {
			if cf.edges[1] != nf.edges[3] {
				continue
			}
		}
		if y > 0 {
			af := pathFrames[(y-1)*sideLen+x]
			if af.edges[2] != nf.edges[0] {
				continue
			}
		}

		out = append(out, nfid)
	}

	return out
}

type tile struct {
	id   string
	data string

	framecache []frame
}

func (t *tile) frames() []frame {
	if t.framecache != nil {
		return t.framecache
	}

	top := t.data[:linelen]
	bottom := t.data[linelen*(linelen-1):]

	var left, right string
	for i := 0; i < linelen; i++ {
		left += string(t.data[linelen*i])
		right += string(t.data[linelen*i+linelen-1])
	}

	edges := []string{
		top,
		right,
		bottom,
		left,
	}

	var frames []frame

	frames = append(frames, frame{t: t, name: "orig", edges: edges, data: func() string { return t.data }})
	frames = append(frames, frame{t: t, name: "orig-r1", edges: []string{reverse(left), top, reverse(right), bottom}, data: func() string { return rotate(t.data) }})
	frames = append(frames, frame{t: t, name: "orig-r2", edges: []string{reverse(bottom), reverse(left), reverse(top), reverse(right)}, data: func() string { return rotate(rotate(t.data)) }})
	frames = append(frames, frame{t: t, name: "orig-r3", edges: []string{right, reverse(bottom), left, reverse(top)}, data: func() string { return rotate(rotate(rotate(t.data))) }})

	top, bottom = reverse(top), reverse(bottom)
	left, right = right, left

	edges = []string{
		top,
		right,
		bottom,
		left,
	}

	data := func() string { return flip(t.data) }
	frames = append(frames, frame{t: t, name: "flip-h", edges: edges, data: data})
	frames = append(frames, frame{t: t, name: "flip-h-r1", edges: []string{reverse(left), top, reverse(right), bottom}, data: func() string { return rotate(data()) }})
	frames = append(frames, frame{t: t, name: "flip-h-r2", edges: []string{reverse(bottom), reverse(left), reverse(top), reverse(right)}, data: func() string { return rotate(rotate(data())) }})
	frames = append(frames, frame{t: t, name: "flip-h-r3", edges: []string{right, reverse(bottom), left, reverse(top)}, data: func() string { return rotate(rotate(rotate(data()))) }})

	t.framecache = frames
	return frames
}

type frame struct {
	t *tile

	name string

	// 0: top, 1: right, 2: bottom, 3: left
	edges []string

	data func() string
}

func (f frame) id() string {
	return f.t.id + ":" + f.name
}

func (f frame) String() string {
	out := f.edges[0] + "\n"
	for i := 1; i < linelen-1; i++ {
		out += string(f.edges[3][i]) + strings.Repeat(" ", linelen-2) + string(f.edges[1][i]) + "\n"
	}
	out += f.edges[2]
	return out
}

func (f frame) DataString() string {
	d := f.data()
	var out string
	for i := 0; i < linelen; i++ {
		if i > 0 {
			out += "\n"
		}
		out += d[i*linelen : (i+1)*linelen]
	}
	return out
}

func reverse(s string) string {
	var out string
	for i := len(s) - 1; i >= 0; i-- {
		out += string(s[i])
	}
	return out
}

// flip horizontally
func flip(data string) string {
	var out string
	for i := 0; i < len(data)/linelen; i++ {
		out += reverse(data[i*linelen : (i+1)*linelen])
	}
	return out
}

func flips(data []string) []string {
	out := make([]string, len(data))
	for i, l := range data {
		out[i] = reverse(l)
	}
	return out
}

// rotate clockwise
func rotate(data string) string {
	// 0 1 2 3
	// 4 5 6 7
	// 8 9 a b
	// c d e f

	// c 8 4 0
	// d 9 5 1
	// e a 6 2
	// f b 7 3

	// 0,0 : 0,3
	// 1,0 : 0,2
	// 2,0 : 0,1
	// 3,0 : 0,0
	// 0,1 : 1,3
	// 1,1 : 1,2
	// 2,1 : 1,1

	var out string
	for x := 0; x < linelen; x++ {
		for y := linelen - 1; y >= 0; y-- {
			out += string(data[y*linelen+x])
		}
	}
	return out
}

// rotate clockwise
func rotates(data []string) []string {
	out := make([]string, len(data))
	for x := 0; x < len(data); x++ {
		for y := len(data[0]) - 1; y >= 0; y-- {
			out[x] += string(data[y][x])
		}
	}
	return out
}

func findSeaMonsters(image []string) (int, int) {
	seaMonster := `
                  #
#    ##    ##    ###
 #  #  #  #  #  #
`

	seaMonsterParts := strings.Split(strings.Trim(seaMonster, "\n"), "\n")
	seaMonsterHashes := strings.Count(seaMonster, "#")

	var found int
	for y := 0; y < len(image)-len(seaMonsterParts); y++ {
	nextX:
		for x := 0; x <= len(image[0])-len(seaMonsterParts[1]); x++ {
			for i, m := range seaMonsterParts[1] {
				if m == '#' && image[y][x+i] != '#' {
					continue nextX
				}
			}
			for i, m := range seaMonsterParts[0] {
				if m == '#' && image[y-1][x+i] != '#' {
					continue nextX
				}
			}
			for i, m := range seaMonsterParts[2] {
				if m == '#' && image[y+1][x+i] != '#' {
					continue nextX
				}
			}
			fmt.Println(image[y-1])
			fmt.Println(image[y])
			fmt.Println(image[y+1])
			fmt.Printf("%s^\n", strings.Repeat(" ", x+1))
			found++
		}
	}
	return found, seaMonsterHashes * found
}
