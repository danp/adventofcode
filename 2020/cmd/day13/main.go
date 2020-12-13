package main

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	depart := scaffold.Int(lines[0])

	var routes []route
	for i, r := range strings.Split(lines[1], ",") {
		if r == "x" {
			continue
		}
		routes = append(routes, route{bus: r, interval: scaffold.Int(r), pos: i})
	}

	{
		minPoint := math.MaxInt64
		var minRoute route
		for _, r := range routes {
			point := 0
			for {
				if point >= depart {
					if point < minPoint {
						minPoint = point
						minRoute = r
					}
					break
				}
				point += r.interval
			}
		}

		fmt.Println(minPoint, minRoute, (minPoint-depart)*minRoute.interval)
	}

	// had to cheat for this because I had no idea about "Chinese remainder theorem."
	// https://rosettacode.org/wiki/Chinese_remainder_theorem#Go
	//
	// this problem sounded very similar to https://adventofcode.com/2019/day/12,
	// The N-Body Problem
	n := make([]*big.Int, 0, len(routes))
	a := make([]*big.Int, 0, len(routes))

	for _, r := range routes {
		n = append(n, big.NewInt(int64(r.interval)))
		a = append(a, big.NewInt(int64((-r.pos)%r.interval)))
	}

	fmt.Println(crt(a, n))
}

type route struct {
	bus      string
	interval int
	pos      int
}

var one = big.NewInt(1)

func crt(a, n []*big.Int) (*big.Int, error) {
	p := new(big.Int).Set(n[0])
	for _, n1 := range n[1:] {
		p.Mul(p, n1)
	}
	var x, q, s, z big.Int
	for i, n1 := range n {
		q.Div(p, n1)
		z.GCD(nil, &s, n1, &q)
		if z.Cmp(one) != 0 {
			return nil, fmt.Errorf("%d not coprime", n1)
		}
		x.Add(&x, s.Mul(a[i], s.Mul(&s, &q)))
	}
	return x.Mod(&x, p), nil
}
