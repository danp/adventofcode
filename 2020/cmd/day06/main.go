package main

import (
	"fmt"
	"os"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	if len(os.Args) > 1 && os.Args[1] == "part2" {
		var sum int
		var ppl int
		data := make(map[rune]int)
		for _, l := range lines {
			if l == "" {
				for _, c := range data {
					if c == ppl {
						sum += 1
					}
				}
				ppl = 0
				data = make(map[rune]int)
				continue
			}
			ppl++
			for _, c := range l {
				data[c]++
			}
		}
		for _, c := range data {
			if c == ppl {
				sum += 1
			}
		}

		fmt.Println(sum)
		return
	}

	var sum int
	data := make(map[rune]bool)
	for _, l := range lines {
		if l == "" {
			sum += len(data)
			data = make(map[rune]bool)
			continue
		}
		for _, c := range l {
			data[c] = true
		}
	}

	sum += len(data)
	fmt.Println(sum)
}
