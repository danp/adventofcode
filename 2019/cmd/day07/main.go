package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/danp/adventofcode/2019/intcode"
	"golang.org/x/sync/errgroup"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	program, err := intcode.Parse(string(b))
	if err != nil {
		panic(err)
	}

	const ampCount = 5

	switch os.Args[1] {
	case "amps":
		mos, mop, err := run(program, ampCount, []int{0, 1, 2, 3, 4}, nil)
		if err != nil {
			panic(err)
		}
		fmt.Println(mos)
		fmt.Println(mop)
	case "feedback":
		mos, mop, err := runFeedback(program, ampCount, []int{5, 6, 7, 8, 9})
		if err != nil {
			panic(err)
		}
		fmt.Println(mos)
		fmt.Println(mop)
	}
}

func run(program []int, ampCount int, phases []int, setup func(amps []*amp, amp *amp, i int)) (int, []int, error) {
	var mos int
	mop := make([]int, ampCount)
	var gErr error

	perm(phases, func(p []int) {
		var g errgroup.Group
		amps := make([]*amp, 0, ampCount)
		for i := 0; i < ampCount; i++ {
			amp := newAmp(i, p[i])
			if i == 0 {
				amp.in <- 0 // first amp signal
			} else {
				amp.in = amps[i-1].out
			}
			if setup != nil {
				setup(amps, amp, i)
			}
			amps = append(amps, amp)
			g.Go(func() error { return amp.run(program) })
		}

		if err := g.Wait(); err != nil {
			gErr = err
			return
		}

		output := <-amps[len(amps)-1].out
		if output > mos {
			mos = output
			copy(mop, p)
		}
	}, 0)

	return mos, mop, gErr
}

func runFeedback(program []int, ampCount int, phases []int) (int, []int, error) {
	return run(program, ampCount, phases, func(amps []*amp, amp *amp, i int) {
		if i == ampCount-1 {
			amp.out = amps[0].in // hook last amp back up to first
		}
	})
}

type amp struct {
	pos     int
	phase   int
	in, out chan int
}

func newAmp(pos, phase int) *amp {
	return &amp{
		pos:   pos,
		phase: phase,
		in:    make(chan int, 1),
		out:   make(chan int, 1),
	}
}

func (a *amp) run(program []int) error {
	inputState := 0
	input := func() (int, error) {
		defer func() { inputState++ }()
		switch inputState {
		case 0:
			return a.phase, nil
		default:
			out := <-a.in
			return out, nil
		}
	}

	output := func(x int) error {
		a.out <- x
		return nil
	}

	p := make([]int, len(program))
	copy(p, program)
	mem := make([]int, len(p))
	return intcode.Run(p, mem, input, output)
}

func perm(a []int, f func([]int), i int) {
	if i > len(a) {
		f(a)
		return
	}
	perm(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
}
