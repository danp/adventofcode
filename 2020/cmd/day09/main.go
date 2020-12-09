package main

import (
	"fmt"
	"sort"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	ints := scaffold.Ints(scaffold.Lines())

	w := window{
		preamb: 25,
	}

	var inv int
	for _, i := range ints {
		if !w.add(i) {
			inv = i
			break
		}
	}

	fmt.Println("invalid number:", inv)

	sums := make(map[int]int)
	for i, n := range ints {
		if n >= inv {
			break
		}
		for j := 0; j <= i; j++ {
			sums[j] += n
			if sums[j] != inv {
				continue
			}

			fmt.Println("sequence found from index", j, "to index", i)
			seq := make([]int, 0)
			ssum := 0
			for k := j; k <= i; k++ {
				seq = append(seq, ints[k])
				ssum += ints[k]
			}
			fmt.Println(seq)
			if ssum == inv {
				fmt.Println("sequence works out")
				sort.Ints(seq)
				min, max := seq[0], seq[len(seq)-1]
				fmt.Println(min, "+", max, "=", min+max)
			} else {
				fmt.Println(ssum, "mismatch", inv-ssum)
			}
			return
		}
	}
}

type window struct {
	preamb int
	cur    []int
}

func (w *window) add(i int) bool {
	if len(w.cur) < w.preamb {
		w.cur = append(w.cur, i)
		return true
	}

	found := false
	for _, o := range w.cur {
		for _, p := range w.cur {
			if o+p == i {
				found = true
				break
			}
		}
	}
	if !found {
		return false
	}

	w.cur = append(w.cur, i)
	w.cur = w.cur[1:]

	return true
}
