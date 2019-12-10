package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFindBestAsteroid(t *testing.T) {
	for i, tt := range []struct {
		grid string
		want result
	}{
		{grid1, result{point{3, 4}, 8}},
		{grid2, result{point{5, 8}, 33}},
		{grid3, result{point{1, 2}, 35}},
		{grid4, result{point{6, 3}, 41}},
		{grid5, result{point{11, 13}, 210}},
	} {
		t.Run(fmt.Sprintf("grid%d", i+1), func(t *testing.T) {
			grid, err := parseGrid(strings.TrimSpace(tt.grid))
			if err != nil {
				t.Fatal(err)
			}

			got, err := findBestAsteroid(grid)
			if err != nil {
				t.Fatal(err)
			}

			if d := cmp.Diff(tt.want, got, cmp.AllowUnexported(result{}, point{})); d != "" {
				t.Errorf("result mismatch (-want +got):\n%s", d)
			}
		})
	}
}

func TestVaporize(t *testing.T) {
	grid, err := parseGrid(strings.TrimSpace(grid5))
	if err != nil {
		t.Fatal(err)
	}

	station := point{11, 13}

	vaporized := vaporize(grid, station)

	for _, tt := range []struct {
		n    int // 1-indexed to make looking at examples easier
		want point
	}{
		{1, point{11, 12}},
		{2, point{12, 1}},
		{3, point{12, 2}},
		{10, point{12, 8}},
		{20, point{16, 0}},
		{50, point{16, 9}},
		{100, point{10, 16}},
		{199, point{9, 6}},
		{200, point{8, 2}},
		{201, point{10, 9}},
		{299, point{11, 1}},
	} {
		got := vaporized[tt.n-1]
		if !got.equal(tt.want) {
			t.Errorf("got point %d %v, want %v", tt.n, got, tt.want)
		}
	}
}

const (
	grid1 = `
.#..#
.....
#####
....#
...##
`

	grid2 = `
......#.#.
#..#.#....
..#######.
.#.#.###..
.#..#.....
..#....#.#
#..#....#.
.##.#..###
##...#..#.
.#....####
`

	grid3 = `
#.#...#.#.
.###....#.
.#....#...
##.#.#.#.#
....#.#.#.
.##..###.#
..#...##..
..##....##
......#...
.####.###.
`

	grid4 = `
.#..#..###
####.###.#
....###.#.
..###.##.#
##.##.#.#.
....###..#
..#.#..#.#
#..#.#.###
.##...##.#
.....#.#..
`

	grid5 = `
.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##
`
)
