package main

import (
	"fmt"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	fmt.Println(find(lines[0], 4))
	fmt.Println(find(lines[0], 14))
}

func find(s string, n int) int {
next:
	for i := 0; i < len(s)-n; i++ {
		p := s[i : i+n]
		cs := make(map[rune]int)
		for _, c := range p {
			if cs[c] > 0 {
				continue next
			}
			cs[c]++
		}
		return i + n
	}
	return -1
}
