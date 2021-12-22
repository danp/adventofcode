package main

import (
	"fmt"
	"regexp"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	// on x=-4..48,y=-30..24,z=-39..15
	var re = regexp.MustCompile(`(\w+) x=([-\d]+)\.\.([-\d]+),y=([-\d]+)\.\.([-\d]+),z=([-\d]+)\.\.([-\d]+)`)

	var cmds []cmd

	for _, l := range lines {
		m := re.FindStringSubmatch(l)

		c := cmd{
			op: m[1],
			cuboid: cuboid{
				min: point3{
					x: scaffold.Int(m[2]),
					y: scaffold.Int(m[4]),
					z: scaffold.Int(m[6]),
				},
				max: point3{
					x: scaffold.Int(m[3]),
					y: scaffold.Int(m[5]),
					z: scaffold.Int(m[7]),
				},
			},
		}
		cmds = append(cmds, c)
	}

	ons := make(map[point3]struct{})
	for _, c := range cmds {
		if c.cuboid.min.x < -50 || c.cuboid.max.x > 50 {
			continue
		}

		switch c.op {
		case "on":
			for x := c.cuboid.min.x; x <= c.cuboid.max.x; x++ {
				for y := c.cuboid.min.y; y <= c.cuboid.max.y; y++ {
					for z := c.cuboid.min.z; z <= c.cuboid.max.z; z++ {
						ons[point3{x, y, z}] = struct{}{}
					}
				}
			}
		case "off":
			for x := c.cuboid.min.x; x <= c.cuboid.max.x; x++ {
				for y := c.cuboid.min.y; y <= c.cuboid.max.y; y++ {
					for z := c.cuboid.min.z; z <= c.cuboid.max.z; z++ {
						delete(ons, point3{x, y, z})
					}
				}
			}
		}
	}

	fmt.Printf("len(ons): %v\n", len(ons))
}

type cmd struct {
	op     string
	cuboid cuboid
}

type point3 struct {
	x, y, z int
}

type cuboid struct {
	min point3
	max point3
}
