package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	reactions, err := parseInput(strings.TrimSpace(string(b)))
	if err != nil {
		panic(err)
	}

	surplus := make(map[string]int)
	fmt.Println(countOre(reactions, "FUEL", 1, surplus))

	x := sort.Search(1000000000000, func(n int) bool {
		return countOre(reactions, "FUEL", n, surplus) > 1000000000000
	}) - 1
	fmt.Println(x)
}

func parseInput(input string) (map[string]reaction, error) {
	out := map[string]reaction{
		"ORE": reaction{q: 1, of: "ORE"},
	}

	lines := strings.Split(input, "\n")
	for _, l := range lines {
		parts := strings.Split(l, "=>")

		name, quant, err := scanPair(parts[1])
		if err != nil {
			return nil, err
		}

		r := reaction{
			q:   quant,
			of:  name,
			req: make(map[string]int),
		}

		parts = strings.Split(parts[0], ",")
		for _, p := range parts {
			name, quant, err := scanPair(p)
			if err != nil {
				return nil, err
			}

			r.req[name] = quant
		}

		out[name] = r
	}

	return out, nil
}

func countOre(reactions map[string]reaction, want string, quantity int, surplus map[string]int) int {
	if want == "ORE" {
		fmt.Println("consuming", quantity, "ORE")
		return quantity
	}

	fmt.Println("want", quantity, want)

	if surp := surplus[want]; surp > 0 {
		fmt.Println("there are", surp, "surplus", want)
		if surp >= quantity {
			surplus[want] -= quantity
			fmt.Println("taking", quantity, want, "from surplus, now", surplus[want], "surplus", want)
			return 0
		}
		quantity -= surplus[want]
		surplus[want] = 0
		fmt.Println("still need", quantity, want, "and surplus now depleted")
		return countOre(reactions, want, quantity, surplus)
	}

	re := reactions[want]
	needed := (quantity-1)/re.q + 1

	fmt.Println("need", needed, "reaction of", want)

	var oreNeeded int
	for rn, rq := range re.req {
		fmt.Println("starting dependent reaction of", rq*needed, rn)
		oreNeeded += countOre(reactions, rn, rq*needed, surplus)
	}

	if excess := needed*re.q - quantity; excess > 0 {
		surplus[want] += excess
		fmt.Println("produced", excess, "excess", want, "and now have surplus of", surplus[want])
	}

	fmt.Println("needed", oreNeeded, "ORE")
	return oreNeeded
}

type reaction struct {
	q  int
	of string

	req map[string]int
}

func scanPair(s string) (string, int, error) {
	var (
		q int
		n string
	)
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%d %s", &q, &n)
	return n, q, err
}
