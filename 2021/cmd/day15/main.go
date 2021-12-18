package main

import (
	"fmt"
	"image"
	"math"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	graph := make(map[image.Point]int)
	var max image.Point
	for y, l := range lines {
		for x, c := range l {
			pt := image.Pt(x, y)
			max = pt
			graph[pt] = int(c) - '0'
		}
	}

	var risk int
	for _, pt := range dijkstra(graph, image.Point{}, max) {
		risk += graph[pt]
	}
	fmt.Println("part 1 risk:", risk)
}

func dijkstra(graph map[image.Point]int, start, end image.Point) []image.Point {
	dists := map[image.Point]int{start: 0}
	prev := make(map[image.Point]image.Point)
	q := make(map[image.Point]struct{})
	for pt := range graph {
		q[pt] = struct{}{}
	}

	for len(q) > 0 {
		pt := min(q, dists)
		delete(q, pt)

		if pt == end {
			break
		}

		for _, npt := range neighbs(q, pt) {
			alt := dists[pt] + graph[npt]
			if nd, ok := dists[npt]; !ok || alt < nd {
				dists[npt] = alt
				prev[npt] = pt
			}
		}
	}

	var path []image.Point
	u := end
	for u != start {
		path = append(path, u)
		u = prev[u]
	}
	return path
}

func min(q map[image.Point]struct{}, dists map[image.Point]int) image.Point {
	mind := math.MaxInt
	var minp image.Point
	for pt := range q {
		if d, ok := dists[pt]; ok && d < mind {
			mind = d
			minp = pt
		}
	}
	if mind == math.MaxInt {
		for pt := range q {
			return pt
		}
	}
	return minp
}

var dirs = []image.Point{
	image.Pt(0, -1),
	image.Pt(1, 0),
	image.Pt(0, 1),
	image.Pt(-1, 0),
}

func neighbs(graph map[image.Point]struct{}, pt image.Point) []image.Point {
	var out []image.Point
	for _, d := range dirs {
		npt := pt.Add(d)
		if _, ok := graph[npt]; ok {
			out = append(out, npt)
		}
	}
	return out
}
