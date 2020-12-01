package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	mode := "part1"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	nums := scaffold.Ints(scaffold.Lines())
	sort.Ints(nums)

	switch mode {
	case "part1":
	outerp1:
		for _, n := range nums {
			for _, o := range nums {
				sum := n + o
				if sum == 2020 {
					fmt.Println(n, o, n*o)
					break outerp1
				}
				if sum > 2020 {
					break
				}
			}
		}
	case "part2":
	outerp2:
		for _, n := range nums {
			for _, o := range nums {
				sum := n + o
				if sum >= 2020 {
					continue
				}

				for _, p := range nums {
					sum := sum + p
					if sum == 2020 {
						fmt.Println(n, o, p, n*o*p)
						break outerp2
					}
					if sum > 2020 {
						break
					}
				}
			}
		}
	}
}
