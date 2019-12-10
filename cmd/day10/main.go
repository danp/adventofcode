package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	grid, err := parseGrid(strings.TrimSpace(string(b)))
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "best":
		res, err := findBestAsteroid(grid)
		if err != nil {
			panic(err)
		}

		fmt.Println(res)
	case "vaporize":
		station := point{20, 18}
		vaporized := vaporize(grid, station)
		th := vaporized[199]
		fmt.Println(th, th.x*100+th.y)
	}
}

type point struct {
	x, y int
}

func (p point) equal(o point) bool {
	return o.x == p.x && o.y == p.y
}

func (p point) dist(o point) float64 {
	return math.Sqrt(math.Pow(float64(o.x-p.x), 2) + math.Pow(float64(o.y-p.y), 2))
}

func (p point) angle(o point) float64 {
	x := math.Atan2(float64(o.y-p.y), float64(o.x-p.x)) * 180 / math.Pi
	if x >= -90 {
		return x + 90
	}
	return 270 + (180 + x)
}

type result struct {
	p point
	c int
}

func findBestAsteroid(grid [][]int) (result, error) {
	points := asteroidPoints(grid)

	var res result
	for _, p := range points {
		c := len(visibleAsteroids(grid, p))
		if c > res.c {
			res.p = p
			res.c = c
		}
	}

	return res, nil

}

func vaporize(grid [][]int, station point) []point {
	var out []point

	for {
		visible := visibleAsteroids(grid, station)
		if len(visible) == 0 {
			break
		}

		sort.Slice(visible, func(i, j int) bool { return visible[i].a < visible[j].a })

		for _, v := range visible {
			out = append(out, v.p)
			grid[v.p.y][v.p.x] = 0
		}
	}
	return out
}

type anglePoint struct {
	p point
	a float64
}

func (a anglePoint) String() string {
	return fmt.Sprintf("(%d,%d) @ %.3f", a.p.x, a.p.y, a.a)
}

func visibleAsteroids(grid [][]int, source point) []anglePoint {
	seen := make(map[float64]bool)
	points := asteroidPoints(grid)
	sort.Slice(points, func(i, j int) bool { return points[i].dist(source) < points[j].dist(source) })
	var out []anglePoint

	for _, p := range points {
		if p.equal(source) {
			continue
		}

		angle := source.angle(p)
		if seen[angle] {
			continue // blocked
		}
		seen[angle] = true
		out = append(out, anglePoint{p, angle})
	}

	return out
}

func asteroidPoints(grid [][]int) []point {
	points := make([]point, 0)
	for y, row := range grid {
		for x, v := range row {
			if v == 1 {
				points = append(points, point{x, y})
			}
		}
	}
	return points
}

func parseGrid(in string) ([][]int, error) {
	lines := strings.Split(in, "\n")
	ll := len(lines[0])

	grid := make([][]int, len(lines))

	for i, l := range lines {
		if tll := len(l); tll != ll {
			return nil, fmt.Errorf("incorrect line length %d, first line length %d", tll, ll)
		}

		grid[i] = make([]int, ll)

		for j, c := range l {
			switch c {
			case '#':
				grid[i][j] = 1
			}
		}
	}

	return grid, nil
}
