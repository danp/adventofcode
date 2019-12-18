package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	maze, err := parse(strings.TrimSpace(string(b)))
	if err != nil {
		panic(err)
	}

	sr := newScorer(maze.g)
	fmt.Println(sr.score(maze.positions, nil))

	ds := newDijkstraScorer(maze.g)
	fmt.Println(ds.score(maze.positions))
}

type skey struct {
	pts  string
	keys string
}

type scorer struct {
	grid  map[image.Point]content
	cache map[skey]int
}

func newScorer(grid map[image.Point]content) *scorer {
	return &scorer{grid: grid, cache: make(map[skey]int)}
}

func (s *scorer) score(pts []image.Point, havekeys map[string]bool) int {
	sk := s.sk(pts, havekeys)
	if c, ok := s.cache[sk]; ok {
		return c
	}

	var kscores []int
	for i, pt := range pts {
		for kpt, kv := range visibleKeys(s.grid, pt, havekeys) {
			hk := mergeKey(havekeys, kv.n)
			newpts := mergePoint(pts, i, kpt)
			ks := s.score(newpts, hk) + kv.s
			kscores = append(kscores, ks)
		}
	}

	if len(kscores) == 0 {
		return 0
	}

	sort.Ints(kscores)
	out := kscores[0]
	s.cache[sk] = out
	return out
}

func (s *scorer) sk(pts []image.Point, havekeys map[string]bool) skey {
	return skey{pts: joinPoints(pts), keys: joinKeys(havekeys)}
}

type dScorer struct {
	grid map[image.Point]content
}

func newDijkstraScorer(grid map[image.Point]content) *dScorer {
	return &dScorer{grid: grid}
}

type dstate struct {
	pts  string
	keys string
}

func (d *dScorer) score(pts []image.Point) int {
	totalkeys := 0
	for _, con := range d.grid {
		if con.key != "" {
			totalkeys++
		}
	}

	pdstate := dstate{pts: joinPoints(pts)}
	dist := map[dstate]int{pdstate: 0}
	prev := make(map[dstate]dstate)
	q := map[dstate]bool{pdstate: true}
	seen := make(map[dstate]bool)

	minDist := func() dstate {
		md := -1
		var mk dstate
		for k, d := range dist {
			if _, ok := q[k]; ok {
				if md == -1 || d < md {
					md = d
					mk = k
				}
			}
		}
		if md == -1 {
			for k := range q {
				return k
			}
		}
		return mk
	}

	for len(q) > 0 {
		u := minDist()
		delete(q, u)
		seen[u] = true

		havekeys := splitKeys(u.keys)
		if len(havekeys) == totalkeys {
			return dist[u]
		}

		upts := splitPoints(u.pts)
		for i, upt := range upts {
			for kpt, kv := range visibleKeys(d.grid, upt, havekeys) {
				np := mergePoint(upts, i, kpt)
				hk := mergeKey(havekeys, kv.n)
				st := dstate{pts: joinPoints(np), keys: joinKeys(hk)}
				if seen[st] {
					continue
				}
				q[st] = true

				alt := dist[u] + kv.s
				if cd, ok := dist[st]; !ok || alt < cd {
					dist[st] = alt
					prev[st] = u
				}
			}
		}
	}

	return 0
}

func joinPoints(pts []image.Point) string {
	spts := make([]image.Point, len(pts))
	for i := 0; i < len(pts); i++ {
		spts[i] = pts[i]
	}
	opts := make([]string, len(spts))
	for i := 0; i < len(spts); i++ {
		opts[i] = fmt.Sprintf("%d,%d", spts[i].X, spts[i].Y)
	}
	return strings.Join(opts, ";")
}

func splitPoints(points string) []image.Point {
	parts := strings.Split(points, ";")
	pts := make([]image.Point, len(parts))
	for i, p := range parts {
		_, err := fmt.Sscanf(p, "%d,%d", &pts[i].X, &pts[i].Y)
		if err != nil {
			panic("unable to split points")
		}
	}
	return pts
}

func mergePoint(pts []image.Point, idx int, newPt image.Point) []image.Point {
	newpts := make([]image.Point, len(pts))
	for i := 0; i < len(pts); i++ {
		newpts[i] = pts[i]
	}
	newpts[idx] = newPt
	return newpts
}

func joinKeys(havekeys map[string]bool) string {
	keys := make([]string, 0, len(havekeys))
	for k := range havekeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, "")
}

func splitKeys(keys string) map[string]bool {
	havekeys := make(map[string]bool)
	for _, c := range keys {
		havekeys[string(c)] = true
	}
	return havekeys
}

func mergeKey(havekeys map[string]bool, newKey string) map[string]bool {
	hk := make(map[string]bool)
	for kn := range havekeys {
		hk[kn] = true
	}
	hk[newKey] = true
	return hk
}

type keyvis struct {
	n string
	s int
}

func visibleKeys(grid map[image.Point]content, pt image.Point, haveKeys map[string]bool) map[image.Point]keyvis {
	out := make(map[image.Point]keyvis)
	seen := make(map[image.Point]bool)

	q := [][]image.Point{{pt}}
	for len(q) > 0 {
		cp := q[0]
		q = q[1:]

		clp := cp[len(cp)-1]
		con := grid[clp]

		if con.door != "" && !haveKeys[strings.ToLower(con.door)] {
			continue
		}

		if con.key != "" && !clp.Eq(pt) && !haveKeys[con.key] {
			out[clp] = keyvis{n: con.key, s: len(cp) - 1}
			continue
		}

		for _, npt := range cards(clp) {
			if _, ok := grid[npt]; ok && !seen[npt] {
				seen[npt] = true

				newp := make([]image.Point, len(cp)+1)
				copy(newp, cp)
				newp[len(newp)-1] = npt
				q = append(q, newp)
			}
		}
	}

	return out
}

func cards(pt image.Point) []image.Point {
	return []image.Point{
		pt.Add(image.Pt(0, 1)),
		pt.Add(image.Pt(1, 0)),
		pt.Sub(image.Pt(0, 1)),
		pt.Sub(image.Pt(1, 0)),
	}
}

type maze struct {
	g         map[image.Point]content
	pos       image.Point
	positions []image.Point
}

type content struct {
	key    string
	door   string
	travel bool
}

func (c content) String() string {
	if c.key != "" {
		return "key(" + c.key + ")"
	} else if c.door != "" {
		return "door(" + c.door + ")"
	} else if c.travel {
		return "travel"
	} else {
		return fmt.Sprintf("%#+v", c)
	}
}

func parse(input string) (*maze, error) {
	g := make(map[image.Point]content)
	var positions []image.Point

	lines := strings.Split(input, "\n")
	for y, l := range lines {
		for x, c := range l {
			pt := image.Pt(x, y)

			var con content
			if c == '#' {
				continue
			} else if c == '.' {
				con.travel = true
			} else if c >= 'a' && c <= 'z' {
				con.key = string(c)
			} else if c >= 'A' && c <= 'Z' {
				con.door = string(c)
			} else if c == '@' {
				positions = append(positions, pt)
				con.travel = true
			} else {
				return nil, fmt.Errorf("unknown char %q", string(c))
			}

			g[pt] = con
		}
	}

	return &maze{g: g, pos: positions[0], positions: positions}, nil
}
