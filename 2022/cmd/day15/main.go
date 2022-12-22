package main

import (
	"fmt"
	"image"
	"math"
	"os"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	part1 := len(os.Args) > 1 && os.Args[0] == "1"

	var bMinX, bMaxX int
	sensors := make(map[image.Point]sensor) // closest beacon
	beacons := make(map[image.Point]struct{})
	for _, l := range lines {
		var s, b image.Point
		if _, err := fmt.Sscanf(l, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &s.X, &s.Y, &b.X, &b.Y); err != nil {
			panic(err)
		}
		sensors[s] = sensor{b, dist(s, b)}
		beacons[b] = struct{}{}
		if bMinX == 0 || b.X < bMinX {
			bMinX = b.X
		}
		if bMaxX == 0 || b.X > bMaxX {
			bMaxX = b.X
		}
	}
	fmt.Printf("sensors: %v\n", sensors)

	if part1 {
		testY := 2000000

		var cannot int
	next1:
		for x := bMinX - 1; ; x-- {
			pt := image.Point{x, testY}
			if _, ok := beacons[pt]; ok {
				continue
			}
			for spt, s := range sensors {
				if d := dist(pt, spt); d <= s.d {
					cannot++
					continue next1
				}
			}
			break
		}
	next2:
		for x := bMaxX + 1; ; x++ {
			pt := image.Point{x, testY}
			if _, ok := beacons[pt]; ok {
				continue
			}
			for spt, s := range sensors {
				if dist(pt, spt) <= s.d {
					cannot++
					continue next2
				}
			}
			break
		}

	next:
		for x := bMinX; x <= bMaxX; x++ {
			pt := image.Point{x, testY}
			if _, ok := beacons[pt]; ok {
				continue
			}
			for spt, s := range sensors {
				if dist(pt, spt) <= s.d {
					cannot++
					continue next
				}
			}
		}
		fmt.Printf("cannot: %v\n", cannot)
		return
	}

	acs := make(map[int]struct{})
	bcs := make(map[int]struct{})
	for spt, s := range sensors {
		acs[spt.Y-spt.X+s.d+1] = struct{}{}
		acs[spt.Y-spt.X-s.d-1] = struct{}{}
		bcs[spt.X+spt.Y+s.d+1] = struct{}{}
		bcs[spt.X+spt.Y-s.d-1] = struct{}{}
	}

	for a := range acs {
	nextb:
		for b := range bcs {
			pt := image.Point{(b - a) / 2, (a + b) / 2}
			if pt.X < 0 || pt.X > 4000000 {
				continue
			}
			if pt.Y < 0 || pt.Y > 4000000 {
				continue
			}
			for spt, s := range sensors {
				if dist(pt, spt) <= s.d {
					continue nextb
				}
			}
			fmt.Printf("pt: %v\n", pt)
			x := pt.X*4000000 + pt.Y
			fmt.Printf("x: %v\n", x)
		}
	}
}

type sensor struct {
	beacon image.Point
	d      int
}

func dist(a, b image.Point) int {
	return int(math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y)))
}
