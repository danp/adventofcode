package main

import (
	"fmt"
	"math"
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

	track(bodies, 1000)
	fmt.Println(energy(bodies))
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

func track(bodies []*body, n int) {
	for iter := 0; iter < n; iter++ {
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
