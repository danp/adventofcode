package main

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDistance(t *testing.T) {
	for _, tt := range []struct {
		line1, line2 string
		distance     int
	}{
		{
			"R75,D30,R83,U83,L12,D49,R71,U7,L72", "U62,R66,U55,R34,D71,R55,D58,R83", 159,
		},
		{
			"R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51", "U98,R91,D20,R16,D67,R40,U7,R15,U6,R7", 135,
		},
	} {
		t.Run(tt.line1+" "+tt.line2, func(t *testing.T) {
			l1, err := parseLine(tt.line1)
			if err != nil {
				t.Fatal(err)
			}

			l2, err := parseLine(tt.line2)
			if err != nil {
				t.Fatal(err)
			}

			if got := minPointDistance(l1.intersections(l2)); got != tt.distance {
				t.Fatalf("got min point distance %d, want %d", got, tt.distance)
			}
		})
	}
}

func TestSteps(t *testing.T) {
	for _, tt := range []struct {
		line1, line2 string
		steps        int
	}{
		{
			"R75,D30,R83,U83,L12,D49,R71,U7,L72", "U62,R66,U55,R34,D71,R55,D58,R83", 610,
		},
		{
			"R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51", "U98,R91,D20,R16,D67,R40,U7,R15,U6,R7", 410,
		},
	} {
		t.Run(tt.line1+" "+tt.line2, func(t *testing.T) {
			l1, err := parseLine(tt.line1)
			if err != nil {
				t.Fatal(err)
			}

			l2, err := parseLine(tt.line2)
			if err != nil {
				t.Fatal(err)
			}

			if got := minPointSteps(l1.pointSteps(), l2.pointSteps()); got != tt.steps {
				t.Fatalf("got min point steps %d, want %d", got, tt.steps)
			}
		})
	}
}

func TestParseLine(t *testing.T) {
	for _, tt := range []struct {
		in   string
		want line
	}{
		{
			"R75,D30,R83,U83,L12,D49,R71,U7,L72",
			line{
				moves: []move{
					{d: right, c: 75},
					{d: down, c: 30},
					{d: right, c: 83},
					{d: up, c: 83},
					{d: left, c: 12},
					{d: down, c: 49},
					{d: right, c: 71},
					{d: up, c: 7},
					{d: left, c: 72},
				},
			},
		},
		{
			"U62,R66,U55,R34,D71,R55,D58,R83",
			line{
				moves: []move{
					{d: up, c: 62},
					{d: right, c: 66},
					{d: up, c: 55},
					{d: right, c: 34},
					{d: down, c: 71},
					{d: right, c: 55},
					{d: down, c: 58},
					{d: right, c: 83},
				},
			},
		},
	} {
		t.Run(tt.in, func(t *testing.T) {
			l, err := parseLine(tt.in)
			if err != nil {
				t.Fatal(err)
			}

			if d := cmp.Diff(tt.want, l, cmp.AllowUnexported(line{}, move{})); d != "" {
				t.Fatalf("line mismatch (-want +got):\n%s", d)
			}
		})
	}

}

func TestLinePoints(t *testing.T) {
	l, err := parseLine("U5,R5,U2,R3,D8,R2,D5,L4")
	if err != nil {
		t.Fatal(err)
	}

	got := l.points()
	want := map[point]bool{
		{y: 1}:         true,
		{y: 2}:         true,
		{y: 3}:         true,
		{y: 4}:         true,
		{y: 5}:         true,
		{x: 1, y: 5}:   true,
		{x: 2, y: 5}:   true,
		{x: 3, y: 5}:   true,
		{x: 4, y: 5}:   true,
		{x: 5, y: 5}:   true,
		{x: 5, y: 6}:   true,
		{x: 5, y: 7}:   true,
		{x: 6, y: -6}:  true,
		{x: 6, y: 7}:   true,
		{x: 7, y: -6}:  true,
		{x: 7, y: 7}:   true,
		{x: 8, y: -6}:  true,
		{x: 8, y: -1}:  true,
		{x: 8}:         true,
		{x: 8, y: 1}:   true,
		{x: 8, y: 2}:   true,
		{x: 8, y: 3}:   true,
		{x: 8, y: 4}:   true,
		{x: 8, y: 5}:   true,
		{x: 8, y: 6}:   true,
		{x: 8, y: 7}:   true,
		{x: 9, y: -6}:  true,
		{x: 9, y: -1}:  true,
		{x: 10, y: -6}: true,
		{x: 10, y: -5}: true,
		{x: 10, y: -4}: true,
		{x: 10, y: -3}: true,
		{x: 10, y: -2}: true,
		{x: 10, y: -1}: true,
	}

	if d := cmp.Diff(want, got, cmp.AllowUnexported(point{})); d != "" {
		t.Fatalf("points mismatch (-want +got):\n%s", d)
	}
}

func TestLinePointSteps(t *testing.T) {
	l, err := parseLine("U5,R5,U2,R3,D8,R2,D5,L4")
	if err != nil {
		t.Fatal(err)
	}

	got := l.pointSteps()
	want := map[point]int{
		{y: 1}:         1,
		{y: 2}:         2,
		{y: 3}:         3,
		{y: 4}:         4,
		{y: 5}:         5,
		{x: 1, y: 5}:   6,
		{x: 2, y: 5}:   7,
		{x: 3, y: 5}:   8,
		{x: 4, y: 5}:   9,
		{x: 5, y: 5}:   10,
		{x: 5, y: 6}:   11,
		{x: 5, y: 7}:   12,
		{x: 6, y: -6}:  34,
		{x: 6, y: 7}:   13,
		{x: 7, y: -6}:  33,
		{x: 7, y: 7}:   14,
		{x: 8, y: -6}:  32,
		{x: 8, y: -1}:  23,
		{x: 8}:         22,
		{x: 8, y: 1}:   21,
		{x: 8, y: 2}:   20,
		{x: 8, y: 3}:   19,
		{x: 8, y: 4}:   18,
		{x: 8, y: 5}:   17,
		{x: 8, y: 6}:   16,
		{x: 8, y: 7}:   15,
		{x: 9, y: -6}:  31,
		{x: 9, y: -1}:  24,
		{x: 10, y: -6}: 30,
		{x: 10, y: -5}: 29,
		{x: 10, y: -4}: 28,
		{x: 10, y: -3}: 27,
		{x: 10, y: -2}: 26,
		{x: 10, y: -1}: 25,
	}

	if d := cmp.Diff(want, got, cmp.AllowUnexported(point{})); d != "" {
		t.Fatalf("pointSteps mismatch (-want +got):\n%s", d)
	}
}

func TestLineIntersections(t *testing.T) {
	l1, err := parseLine("R75,D30,R83,U83,L12,D49,R71,U7,L72")
	if err != nil {
		t.Fatal(err)
	}

	l2, err := parseLine("U62,R66,U55,R34,D71,R55,D58,R83")
	if err != nil {
		t.Fatal(err)
	}

	got := l1.intersections(l2)
	sort.Slice(got, func(i, j int) bool {
		if got[i].x != got[j].x {
			return got[i].x < got[j].x
		}
		return got[i].y < got[j].y
	})

	want := []point{
		{x: 146, y: 46},
		{x: 155, y: 4},
		{x: 155, y: 11},
		{x: 158, y: -12},
	}

	if d := cmp.Diff(want, got, cmp.AllowUnexported(point{})); d != "" {
		t.Fatalf("points mismatch (-want +got):\n%s", d)
	}
}

func TestMinPointDistance(t *testing.T) {
	points := []point{
		{x: 158, y: -12},
		{x: 146, y: 46},
		{x: 155, y: 4},
		{x: 155, y: 11},
	}

	got := minPointDistance(points)
	if want := 159; got != want {
		t.Fatalf("got min point distance %d, want %d", got, want)
	}
}
