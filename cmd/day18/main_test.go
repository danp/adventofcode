package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestScore(t *testing.T) {
	for _, mi := range []int{8, 86, 132, 136, 81} {
		t.Run(fmt.Sprintf("maze%d", mi), func(t *testing.T) {
			maze, err := parse(strings.TrimSpace(mazes[mi]))
			if err != nil {
				t.Fatal(err)
			}

			sr := newScorer(maze.g)
			if got := sr.score(maze.positions, nil); got != mi {
				t.Errorf("got score %d, want %d", got, mi)
			}
		})
	}
}

func TestDijkstraScore(t *testing.T) {
	for _, mi := range []int{8, 86, 132, 136, 81} {
		t.Run(fmt.Sprintf("maze%d", mi), func(t *testing.T) {
			maze, err := parse(strings.TrimSpace(mazes[mi]))
			if err != nil {
				t.Fatal(err)
			}

			sr := newDijkstraScorer(maze.g)
			if got := sr.score(maze.positions); got != mi {
				t.Errorf("got score %d, want %d", got, mi)
			}
		})
	}
}

func TestMultiScore(t *testing.T) {
	for _, mi := range []int{8, 24} {
		t.Run(fmt.Sprintf("maze%d", mi), func(t *testing.T) {
			maze, err := parse(strings.TrimSpace(multiMazes[mi]))
			if err != nil {
				t.Fatal(err)
			}

			sr := newScorer(maze.g)
			if got := sr.score(maze.positions, nil); got != mi {
				t.Errorf("got score %d, want %d", got, mi)
			}
		})
	}
}

func TestDijkstraMultiScore(t *testing.T) {
	for _, mi := range []int{8, 24} {
		t.Run(fmt.Sprintf("maze%d", mi), func(t *testing.T) {
			maze, err := parse(strings.TrimSpace(multiMazes[mi]))
			if err != nil {
				t.Fatal(err)
			}

			sr := newDijkstraScorer(maze.g)
			if got := sr.score(maze.positions); got != mi {
				t.Errorf("got score %d, want %d", got, mi)
			}
		})
	}
}

var mazes = map[int]string{
	8: `
#########
#b.A.@.a#
#########
`,

	86: `
########################
#f.D.E.e.C.b.A.@.a.B.c.#
######################.#
#d.....................#
########################
`,

	132: `
########################
#...............b.C.D.f#
#.######################
#.....@.a.B.c.d.A.e.F.g#
########################
`,

	136: `
#################
#i.G..c...e..H.p#
########.########
#j.A..b...f..D.o#
########@########
#k.E..a...g..B.n#
########.########
#l.F..d...h..C.m#
#################
`,

	81: `
########################
#@..............ac.GI.b#
###d#e#f################
###A#B#C################
###g#h#i################
########################
`,
}

var multiMazes = map[int]string{
	8: `
#######
#a.#Cd#
##@#@##
#######
##@#@##
#cB#Ab#
#######
`,

	24: `
###############
#d.ABC.#.....a#
######@#@######
###############
######@#@######
#b.....#.....c#
###############
`,
}
