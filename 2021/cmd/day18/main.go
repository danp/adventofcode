package main

import (
	"fmt"
)

func main() {
	// lines := scaffold.Lines()

	r1 := pair{pair{pair{pair{pair{9, 8}, 1}, 2}, 3}, 4}
	r2 := pair{pair{pair{pair{0, 9}, 2}, 3}, 4}

	fmt.Printf("r1: %v\n", r1)
	fmt.Printf("r2: %v\n", r2)
}

type node interface{}

type pair struct {
	a, b node
}

func (p pair) add(o pair) pair {
	return pair{p, o}
}
