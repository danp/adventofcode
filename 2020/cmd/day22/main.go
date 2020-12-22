package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	origDecks := [][]int{{}, {}}
	var st int
	for _, l := range lines {
		if strings.HasPrefix(l, "Player") {
			continue
		}
		if l == "" {
			st++
			continue
		}

		origDecks[st] = append(origDecks[st], scaffold.Int(l))
	}

	{
		decks := make([][]int, len(origDecks))
		copy(decks, origDecks)
		play(-1, decks)

		for i, d := range decks {
			var score int
			for i, c := range d {
				score += (len(d) - i) * c
			}
			fmt.Println("player", i+1, "normal score", score)
		}
	}

	decks := make([][]int, len(origDecks))
	copy(decks, origDecks)
	play(0, decks)

	for i, d := range decks {
		var score int
		for i, c := range d {
			score += (len(d) - i) * c
		}
		fmt.Println("player", i+1, "recursive score", score)
	}
}

func play(depth int, decks [][]int) int {
	fmt.Println(depth, decks)

	seen := map[string]bool{
		deckstr(decks[0]): true,
		deckstr(decks[1]): true,
	}

	for len(decks[0]) > 0 && len(decks[1]) > 0 {
		c0 := decks[0][0]
		c1 := decks[1][0]

		var winner int
		if depth >= 0 && c0 < len(decks[0]) && c1 < len(decks[1]) {
			p0 := make([]int, c0)
			copy(p0, decks[0][1:c0+1])
			p1 := make([]int, c1)
			copy(p1, decks[1][1:c1+1])
			winner = play(depth+1, [][]int{p0, p1})
		} else if c1 > c0 {
			winner = 1
		}

		if winner == 0 {
			decks[0] = append(decks[0], c0, c1)
		} else {
			decks[1] = append(decks[1], c1, c0)
		}
		decks[0] = decks[0][1:]
		decks[1] = decks[1][1:]

		fmt.Println(depth, decks)

		ds0 := deckstr(decks[0])
		ds1 := deckstr(decks[1])
		if d0s, d1s := seen[ds0], seen[ds1]; d0s || d1s {
			fmt.Println(depth, "seen winner 0", d0s, d1s, "ds0", ds0, "ds1", ds1)
			return 0
		}
		seen[ds0] = true
		seen[ds1] = true
	}

	if len(decks[0]) > 0 {
		fmt.Println(depth, "winner 0")
		return 0
	}
	fmt.Println(depth, "winner 1")
	return 1
}

func deckstr(deck []int) string {
	var out string
	for i, c := range deck {
		if i > 0 {
			out += " "
		}
		out += strconv.Itoa(c)
	}
	return out
}
