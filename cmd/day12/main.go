package main

import (
	"fmt"
	"math"
	"os"
)

func main() {
	/*
		<x=6, y=-2, z=-7>
		<x=-6, y=-7, z=-4>
		<x=-9, y=11, z=0>
		<x=-3, y=-4, z=6>
	*/

	bodies := []*body{
		{p: p3{6, -2, -7}},
		{p: p3{-6, -7, -4}},
		{p: p3{-9, 11, 0}},
		{p: p3{-3, -4, 6}},
	}

	if os.Args[1] == "part1" {
		track(bodies, 1000)
		fmt.Println(energy(bodies))
		return
	}

	if os.Args[1] == "part2" {
		seenX := make(map[string]bool)
		seenY := make(map[string]bool)
		seenZ := make(map[string]bool)

		for {
			step(bodies)

			// could comparing to some initial state also work?
			xs := fmt.Sprintf("%v", bmap(bodies, func(b *body) []int { return []int{b.p.x, b.v.x} }))
			ys := fmt.Sprintf("%v", bmap(bodies, func(b *body) []int { return []int{b.p.y, b.v.y} }))
			zs := fmt.Sprintf("%v", bmap(bodies, func(b *body) []int { return []int{b.p.z, b.v.z} }))

			if seenX[xs] && seenY[ys] && seenZ[zs] {
				break
			}

			seenX[xs] = true
			seenY[ys] = true
			seenZ[zs] = true
		}

		fmt.Println(len(seenX), len(seenY), len(seenZ))
		fmt.Println(lcm(len(seenX), len(seenY), len(seenZ)))
		return
	}
}

type p3 struct {
	x, y, z int
}

func (p p3) add(o p3) p3 {
	return p3{x: p.x + o.x, y: p.y + o.y, z: p.z + o.z}
}

type body struct {
	p p3
	v p3
}

func (b body) potential() int {
	return int(math.Abs(float64(b.p.x)) + math.Abs(float64(b.p.y)) + math.Abs(float64(b.p.z)))
}

func (b body) kinetic() int {
	return int(math.Abs(float64(b.v.x)) + math.Abs(float64(b.v.y)) + math.Abs(float64(b.v.z)))
}

func (b body) total() int {
	return b.potential() * b.kinetic()
}

func (b body) copy() body {
	return body{p: b.p, v: b.v}
}

func track(bodies []*body, n int) {
	for iter := 0; iter < n; iter++ {
		step(bodies)
	}
}

func step(bodies []*body) {
	// apply gravity
	for i, c := range bodies {
		for j, o := range bodies {
			if i == j {
				continue
			}

			c.v.x += vchange(c.p.x, o.p.x)
			c.v.y += vchange(c.p.y, o.p.y)
			c.v.z += vchange(c.p.z, o.p.z)
		}
	}

	// apply velocity
	for _, b := range bodies {
		b.p = b.p.add(b.v)
	}
}

func energy(bodies []*body) int {
	var out int
	for _, b := range bodies {
		out += b.total()
	}
	return out
}

func vchange(n1, n2 int) int {
	if n1 == n2 {
		return 0
	}
	if n1 > n2 {
		return -1
	}
	return 1
}

func bmap(bodies []*body, f func(b *body) []int) []int {
	var out []int
	for _, b := range bodies {
		out = append(out, f(b)...)
	}
	return out
}

// greatest common divisor (GCD) via Euclidean algorithm
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func lcm(a, b int, integers ...int) int {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}
