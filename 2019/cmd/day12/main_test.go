package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTracking(t *testing.T) {
	bodies := []*body{
		{p: p3{-1, 0, 2}},
		{p: p3{2, -10, -7}},
		{p: p3{4, -8, 8}},
		{p: p3{3, 5, -1}},
	}

	track(bodies, 1)

	want := []*body{
		{p: p3{2, -1, 1}, v: p3{3, -1, -1}},
		{p: p3{3, -7, -4}, v: p3{1, 3, 3}},
		{p: p3{1, -7, 5}, v: p3{-3, 1, -3}},
		{p: p3{2, 2, 0}, v: p3{-1, -3, 1}},
	}
	if d := cmp.Diff(want, bodies, cmp.AllowUnexported(body{}, p3{})); d != "" {
		t.Errorf("tracking mismatch (-want +got):\n%s", d)
	}
}

func TestTracking10(t *testing.T) {
	bodies := []*body{
		{p: p3{-1, 0, 2}},
		{p: p3{2, -10, -7}},
		{p: p3{4, -8, 8}},
		{p: p3{3, 5, -1}},
	}

	track(bodies, 10)

	want := []*body{
		{p: p3{2, 1, -3}, v: p3{-3, -2, 1}},
		{p: p3{1, -8, 0}, v: p3{-1, 1, 3}},
		{p: p3{3, -6, 1}, v: p3{3, 2, -3}},
		{p: p3{2, 0, 4}, v: p3{1, -1, -1}},
	}
	if d := cmp.Diff(want, bodies, cmp.AllowUnexported(body{}, p3{})); d != "" {
		t.Errorf("tracking mismatch (-want +got):\n%s", d)
	}

	if got, want := energy(bodies), 179; got != want {
		t.Errorf("got energy %d, want %d", got, want)
	}
}
