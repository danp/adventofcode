package main

import (
	"fmt"
	"image"
	"math"

	"github.com/danp/adventofcode/scaffold"
	"github.com/dominikbraun/graph"
)

func main() {
	lines := scaffold.Lines()

	var rect image.Rectangle
	var start, end image.Point
	var as []image.Point
	g := graph.New(func(s square) image.Point { return s.pt }, graph.Directed())
	for y, l := range lines {
		for x, c := range l {
			pt := image.Point{x, y}

			switch c {
			case 'S':
				start = pt
				c = 'a'
			case 'E':
				end = pt
				c = 'z'
			}
			if err := g.AddVertex(square{pt, c}); err != nil {
				panic(err)
			}

			if c == 'a' {
				as = append(as, pt)
			}

			rect.Max = pt
		}
	}

	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		for x := rect.Min.X; x <= rect.Max.X; x++ {
			pt := image.Point{x, y}
			for _, npt := range []image.Point{pt.Add(image.Point{-1, 0}), pt.Add(image.Point{1, 0}), pt.Add(image.Point{0, -1}), pt.Add(image.Point{0, 1})} {
				if npt.X < rect.Min.X || npt.X > rect.Max.X {
					continue
				}
				if npt.Y < rect.Min.Y || npt.Y > rect.Max.Y {
					continue
				}
				if npt == pt {
					continue
				}
				s, err := g.Vertex(pt)
				if err != nil {
					panic(err)
				}
				ns, err := g.Vertex(npt)
				if err != nil {
					panic(err)
				}
				if ns.el > s.el+1 {
					continue
				}
				// weight needed, https://github.com/dominikbraun/graph/issues/70
				if err := g.AddEdge(pt, npt, graph.EdgeWeight(1)); err != nil {
					panic(err)
				}
			}
		}
	}

	fmt.Printf("start: %v\n", start)
	fmt.Printf("end: %v\n", end)
	fmt.Printf("g.Size(): %v\n", g.Size())
	fmt.Printf("g.Order(): %v\n", g.Order())

	ps, err := graph.ShortestPath(g, start, end)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ps: %v\n", ps)
	fmt.Printf("len(ps): %v\n", len(ps)-1)

	// this is very slow
	min := math.MaxInt
	for _, a := range as {
		ps, err := graph.ShortestPath(g, a, end)
		if err != nil {
			continue
		}
		l := len(ps) - 1
		if l < min {
			min = l
		}
	}
	fmt.Printf("min: %v\n", min)
}

type square struct {
	pt image.Point
	el rune
}
