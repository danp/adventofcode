package main

import (
	"fmt"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	var sum int
	for _, l := range lines {
		a, b := runeSet(l[:len(l)/2]), runeSet(l[len(l)/2:])
		for c := range setIntersect(a, b) {
			var pri int
			if c >= 'a' && c <= 'z' {
				pri = int(c) - 'a' + 1
			} else if c >= 'A' && c <= 'Z' {
				pri = int(c) - 'A' + 27
			}
			sum += pri
		}
	}
	fmt.Println(sum)

	sum = 0
	for i := 0; i < len(lines); i += 3 {
		a, b, c := runeSet(lines[i]), runeSet(lines[i+1]), runeSet(lines[i+2])
		for c := range setIntersect(a, b, c) {
			var pri int
			if c >= 'a' && c <= 'z' {
				pri = int(c) - 'a' + 1
			} else if c >= 'A' && c <= 'Z' {
				pri = int(c) - 'A' + 27
			}
			sum += pri
		}

	}
	fmt.Println(sum)
}

func runeSet(s string) map[rune]struct{} {
	set := make(map[rune]struct{})
	for _, r := range s {
		set[r] = struct{}{}
	}
	return set
}

type set map[rune]struct{}

func setIntersect(a set, sets ...set) map[rune]struct{} {
	set := make(map[rune]struct{})
A:
	for r := range a {
		for _, s := range sets {
			if _, ok := s[r]; !ok {
				continue A
			}
		}
		set[r] = struct{}{}
	}
	return set
}
