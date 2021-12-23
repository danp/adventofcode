package main

import (
	"fmt"
	"regexp"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	var ons []cuboid
	for _, l := range lines {
		c := newCmd(l)

		fmt.Printf("c.cuboid.size(): %v\n", c.cuboid.size())

		for _, on := range ons {
			// this won't work because we don't know who's using up what portions

		}
		ons = append(ons, c.cuboid)

		// how many on in given cuboid, how many does cuboid turn on

		// do we need N layers per cuboid

		// on X -> off Y -> on Z

		// if on
		//   for each existing

		// if on
		//  for each existing overlapping
		//  trim overlap, add to that one's on count

		//  if any left, add new existing
		// if off
		//  for each existing overlapping
		//  trim overlap, sub from that one's on count
		//  if any left, whatever
	}

	// sum up on count

}

type cmd struct {
	op     string
	cuboid cuboid
}

var cmdRe = regexp.MustCompile(`(\w+) x=([-\d]+)\.\.([-\d]+),y=([-\d]+)\.\.([-\d]+),z=([-\d]+)\.\.([-\d]+)`)

func newCmd(s string) cmd {
	// on x=-4..48,y=-30..24,z=-39..15
	m := cmdRe.FindStringSubmatch(s)
	return cmd{
		op: m[1],
		cuboid: cuboid{
			min: point3{
				x: scaffold.Int(m[2]),
				y: scaffold.Int(m[4]),
				z: scaffold.Int(m[6]),
			},
			max: point3{
				x: scaffold.Int(m[3]) + 1,
				y: scaffold.Int(m[5]) + 1,
				z: scaffold.Int(m[7]) + 1,
			},
		},
	}
}

type point3 struct {
	x, y, z int
}

type cuboid struct {
	min point3
	max point3
}

func (c cuboid) empty() bool {
	return c.min.x >= c.max.x || c.min.y >= c.max.y || c.min.z >= c.max.z
}

func (c cuboid) size() int {
	return (c.max.x - c.min.x) * (c.max.y - c.min.y) * (c.max.z - c.min.z)
}

func (c cuboid) overlap(o cuboid) cuboid {
	has := c.min.x < o.max.x && o.min.x < c.max.x &&
		c.min.y < o.max.y && o.min.y < c.max.y &&
		c.min.z < o.max.z && o.min.z < c.max.z
	if !has {
		return cuboid{}
	}

	return cuboid{
		min: max(c.min, o.min),
		max: min(c.max, o.max),
	}
}

func min(a, b point3) point3 {
	if a.x < b.x || a.y < b.y || a.z < b.z {
		return a
	}
	return b
}

func max(a, b point3) point3 {
	if a.x > b.x || a.y > b.y || a.z > b.z {
		return a
	}
	return b
}
