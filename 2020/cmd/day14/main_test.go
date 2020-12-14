package main

import "testing"

func TestBitperms(t *testing.T) {
	perms := bitperms([]int{16, 8, 2})
	t.Log(perms)
	if len(perms) != 8 {
		t.Errorf("not enough %d", len(perms))
	}

	vals := make(map[int][]int)
	for _, p := range perms {
		for b, v := range p {
			vals[b] = append(vals[b], v)
		}
	}

	t.Log(vals)

	perms = bitperms([]int{16})
	t.Log(perms)
	if len(perms) != 2 {
		t.Errorf("not enough %d", len(perms))
	}

	vals = make(map[int][]int)
	for _, p := range perms {
		for b, v := range p {
			vals[b] = append(vals[b], v)
		}
	}

	t.Log(vals)
}
