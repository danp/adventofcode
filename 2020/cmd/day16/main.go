package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	var rules []rule
	for _, l := range lines {
		if l == "" {
			break
		}

		nameAndRanges := strings.Split(l, ": ")

		var r rule
		r.name = nameAndRanges[0]

		for _, rng := range strings.Split(nameAndRanges[1], " or ") {
			var min, max int
			if _, err := fmt.Sscanf(rng, "%d-%d", &min, &max); err != nil {
				panic(err)
			}
			r.ranges = append(r.ranges, []int{min, max})
		}

		rules = append(rules, r)
	}

	var myTicket ticket
	var nearbyTickets []ticket
	var invalidValues []int

	var state int
	for _, l := range lines {
		if state == 0 {
			if l == "" {
				state++
			}
			continue
		}

		if state == 1 {
			if l != "your ticket:" {
				myTicket = parseTicket(l)
				state++
				continue
			}
			continue
		}

		if l == "nearby tickets:" {
			continue
		}

		ticket := parseTicket(l)
		valid := true
		for _, v := range ticket.values {
			var oneRule bool
			for _, r := range rules {
				if r.match(v) {
					oneRule = true
					break
				}
			}
			if !oneRule {
				valid = false
				invalidValues = append(invalidValues, v)
			}
		}
		if !valid {
			continue
		}

		nearbyTickets = append(nearbyTickets, ticket)

	}

	var sum int
	for _, v := range invalidValues {
		sum += v
	}
	fmt.Println(invalidValues, sum)

	// val idx -> rule idx -> bool
	poss := make(map[int]map[int]bool)
	for i := range myTicket.values {
		poss[i] = make(map[int]bool)
		for j := range rules {
			poss[i][j] = true
		}
	}

	for _, t := range nearbyTickets {
		for i, v := range t.values {
			for j, r := range rules {
				if _, ok := poss[i][j]; !ok {
					continue
				}

				if !r.match(v) {
					fmt.Println(i, "can't be", r.name)
					delete(poss[i], j)
					continue
				}
			}
		}
	}

	done := make(map[int]int)
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

			var val int
			for j := range rm {
				val = j
				break
			}

			fmt.Println(i, "must be", rules[val].name)

			for j, jrm := range poss {
				if j == i {
					continue
				}

				delete(jrm, val)
			}

			done[i] = val
		}
	}

	// for multiplying
	departureValues := 1
	for vi, ri := range done {
		r := rules[ri]
		fmt.Println(vi, "is", r.name, "which on my ticket is", myTicket.values[vi])

		if strings.HasPrefix(r.name, "departure") {
			departureValues *= myTicket.values[vi]
		}
	}

	fmt.Println(departureValues)
}

type rule struct {
	name   string
	ranges [][]int
}

func (r rule) match(v int) bool {
	for _, rr := range r.ranges {
		if v >= rr[0] && v <= rr[1] {
			return true
		}
	}
	return false
}

type ticket struct {
	values []int
}

func parseTicket(l string) ticket {
	valss := strings.Split(l, ",")
	vals := make([]int, 0, len(valss))
	for _, vs := range valss {
		v, _ := strconv.Atoi(vs)
		vals = append(vals, v)
	}
	return ticket{values: vals}
}
