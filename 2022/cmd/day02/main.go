package main

import (
	"fmt"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	part1(lines)
	part2(lines)
}

func part1(lines []string) {
	values := map[string]int{
		"A": 1, // rock
		"B": 2, // paper
		"C": 3, // scissors
		"X": 1, // rock
		"Y": 2, // paper
		"Z": 3, // scissors
	}

	var score int
	for _, l := range lines {
		opp, me, _ := strings.Cut(l, " ")

		var round int
		switch fmt.Sprint(values[me], values[opp]) {
		case "1 3":
			round = 6
		case "3 2":
			round = 6
		case "2 1":
			round = 6
		default:
			if values[me] == values[opp] {
				round = 3
			}
		}

		round += values[me]

		score += round
	}
	fmt.Println(score)
}

func part2(lines []string) {
	moves := map[string][]string{ // lose, win
		"A": {"Z", "Y"},
		"B": {"X", "Z"},
		"C": {"Y", "X"},
	}

	var out []string
	for _, l := range lines {
		opp, want, _ := strings.Cut(l, " ")

		var m string
		switch want {
		case "X":
			m = moves[opp][0]
		case "Y":
			m = opp
		case "Z":
			m = moves[opp][1]
		}

		out = append(out, opp+" "+m)
	}

	part1(out)
}
