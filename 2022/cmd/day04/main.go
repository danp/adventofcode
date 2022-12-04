package main

import (
	"fmt"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	var contains int
	var overlaps int
	for _, l := range lines {
		a, b, _ := strings.Cut(l, ",")

		ar, br := newRng(a), newRng(b)
		if ar.contains(br) || br.contains(ar) {
			contains++
		}
		if ar.overlaps(br) {
			overlaps++
		}
	}
	fmt.Println(contains)
	fmt.Println(overlaps)
}

type rng struct {
	begin, end int
}

func newRng(s string) rng {
	b, e, _ := strings.Cut(s, "-")
	return rng{scaffold.Int(b), scaffold.Int(e)}
}

func (r rng) contains(o rng) bool {
	return o.begin >= r.begin && o.end <= r.end
}

func (r rng) overlaps(o rng) bool {
	if o.begin > r.end {
		return false
	}
	if o.end < r.begin {
		return false
	}
	return true
}
