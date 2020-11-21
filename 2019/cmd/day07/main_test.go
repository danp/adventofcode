package main

import (
	"testing"

	"github.com/danp/adventofcode/2019/intcode"
	"github.com/google/go-cmp/cmp"
)

func TestAmps(t *testing.T) {
	const prog = "3,31,3,32,1002,32,10,32,1001,31,-2,31,1007,31,0,33,1002,33,7,33,1,33,31,31,1,32,31,31,4,31,99,0,0,0"

	program, err := intcode.Parse(prog)
	if err != nil {
		t.Fatal(err)
	}

	mos, mop, err := run(program, 5, []int{0, 1, 2, 3, 4}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := mos, 65210; got != want {
		t.Errorf("got max output signal %d, want %d", got, want)
	}

	wantPhases := []int{1, 0, 4, 3, 2}
	if d := cmp.Diff(wantPhases, mop); d != "" {
		t.Errorf("max output phases mismatch (-want +got):\n%s", d)
	}
}

func TestFeedback(t *testing.T) {
	const prog = "3,52,1001,52,-5,52,3,53,1,52,56,54,1007,54,5,55,1005,55,26,1001,54,-5,54,1105,1,12,1,53,54,53,1008,54,0,55,1001,55,1,55,2,53,55,53,4,53,1001,56,-1,56,1005,56,6,99,0,0,0,0,10"

	program, err := intcode.Parse(prog)
	if err != nil {
		t.Fatal(err)
	}

	mos, mop, err := runFeedback(program, 5, []int{5, 6, 7, 8, 9})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := mos, 18216; got != want {
		t.Errorf("got max output signal %d, want %d", got, want)
	}

	wantPhases := []int{9, 7, 8, 5, 6}
	if d := cmp.Diff(wantPhases, mop); d != "" {
		t.Errorf("max output phases mismatch (-want +got):\n%s", d)
	}
}
