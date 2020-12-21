package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	ingrs := make(map[string]map[string]int)
	ingrmentions := make(map[string]int)
	allers := make(map[string]int)

	for _, l := range lines {
		var st int
		var lingrs []string
		for _, f := range strings.Fields(l) {
			if f == "(contains" {
				st++
				continue
			}
			if st == 1 {
				aller := strings.TrimSuffix(strings.TrimSuffix(f, ","), ")")
				allers[aller]++

				for _, i := range lingrs {
					if ingrs[i] == nil {
						ingrs[i] = make(map[string]int)
					}
					ingrs[i][aller]++
				}
				continue
			}

			ingrmentions[f]++
			lingrs = append(lingrs, f)
		}
	}

	var times int
outer:
	for i, im := range ingrs {
		for a, am := range allers {
			if im[a] == am {
				continue outer
			}
		}
		times += ingrmentions[i]
	}

	fmt.Println(times, "appearances for inert ingredients")

	poss := make(map[string]map[string]bool)

	for i, im := range ingrs {
		for a, am := range allers {
			if im[a] == am {
				if poss[i] == nil {
					poss[i] = make(map[string]bool)
				}
				poss[i][a] = true
			}
		}
	}

	done := make(map[string]string)
	for {
		if len(done) == len(poss) {
			break
		}

		for i, rm := range poss {
			if len(rm) > 1 {
				continue
			}

			if _, ok := done[i]; ok {
				continue
			}

			var val string
			for j := range rm {
				val = j
				break
			}

			fmt.Println(i, "must be", val)

			for j, jrm := range poss {
				if j == i {
					continue
				}

				delete(jrm, val)
			}

			done[i] = val
		}
	}

	out := make([]string, 0, len(done))
	for i := range done {
		out = append(out, i)
	}
	sort.Slice(out, func(i, j int) bool { return done[out[i]] < done[out[j]] })
	fmt.Println(strings.Join(out, ","))
}
