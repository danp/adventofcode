package main

import (
	"fmt"
)

func main() {
	// lines := scaffold.Lines()

	n1 := pair{1, 2}
	n2 := pair{pair{3, 4}, 5}

	n3 := n1.add(n2)
	fmt.Printf("n3: %v\n", n3)
}

type node interface{}

type pair struct {
	a, b node
}

func (p pair) add(o pair) pair {
	return pair{p, o}
}
