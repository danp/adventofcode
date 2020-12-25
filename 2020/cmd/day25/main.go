package main

import (
	"fmt"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	pubkeys := scaffold.Ints(scaffold.Lines())

	loopSizes := make(map[int]int)
	for _, p := range pubkeys {
		loopSizes[p] = determineLoop(p, 7)
	}

	fmt.Println(loopSizes)
	fmt.Println(transform(pubkeys[0], loopSizes[pubkeys[1]]))
}

func determineLoop(pubkey, subject int) int {
	value := 1

	loop := 0
	for loop < 20201227 {
		loop++

		value *= subject
		value %= 20201227

		if value == pubkey {
			return loop
		}
	}
	panic("didn't get it")
}

func transform(subject, loop int) int {
	value := 1
	for i := 0; i < loop; i++ {
		value *= subject
		value %= 20201227
	}
	return value
}
