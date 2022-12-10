package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(b), "\n")

	var stacks [][]string
	for _, l := range lines {
		if l[:3] == " 1 " {
			break
		}
		for i := 0; ; i++ {
			j := i * 4 // 3 chars for [X] + space
			if len(l) < j+3 {
				break
			}
			p := strings.TrimSpace(l[j : j+3])
			if len(stacks) < i+1 {
				stacks = append(stacks, []string{})
			}
			if p == "" {
				continue
			}
			stacks[i] = append(stacks[i], p)
		}
	}

	for i := 0; i < len(stacks); i++ {
		reverse(stacks[i])
	}

	run(lines, clone(stacks), false)
	run(lines, clone(stacks), true)
}

func run(lines []string, stacks [][]string, part2 bool) {
	for _, l := range lines {
		if !strings.HasPrefix(l, "move") {
			continue
		}
		var q, from, to int
		if _, err := fmt.Sscanf(l, "move %d from %d to %d", &q, &from, &to); err != nil {
			panic(err)
		}
		fs := stacks[from-1]
		moving := fs[len(fs)-q:]
		if !part2 {
			reverse(moving)
		}
		stacks[from-1] = stacks[from-1][:len(stacks[from-1])-q]
		stacks[to-1] = append(stacks[to-1], moving...)
	}

	var lasts string
	for i := 0; i < len(stacks); i++ {
		s := stacks[i]
		lasts += s[len(s)-1]
	}
	lasts = strings.ReplaceAll(lasts, "[", "")
	lasts = strings.ReplaceAll(lasts, "]", "")
	fmt.Println(lasts)
}

func reverse(s []string) {
	for j := len(s)/2 - 1; j >= 0; j-- {
		opp := len(s) - 1 - j
		s[j], s[opp] = s[opp], s[j]
	}
}

func clone(s [][]string) [][]string {
	var new [][]string
	for _, x := range s {
		new = append(new, slices.Clone(x))
	}
	return new
}
