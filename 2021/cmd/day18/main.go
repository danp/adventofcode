package main

import (
	"fmt"
)

type atomInt int

type atomPair struct {
	a, b atomElement
}

func (p atomPair) add(b atomPair) atomPair {
	return atomPair{p, b}
}

type atomElement interface{}

func main() {
	// lines := scaffold.Lines()

	a := atomPair{1, 2}
	b := atomPair{atomPair{3, 4}, 5}

	c := a.add(b)
	fmt.Println(c)

	// [1,2] + [[3,4],5] becomes [[1,2],[[3,4],5]]
}
