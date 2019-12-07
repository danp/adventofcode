package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/danp/adventofcode2019/intcode"
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
	var maxOutputSignal int
	maxPhaseSettings := make([]int, ampCount)
	perm([]int{0, 1, 2, 3, 4}, func(p []int) {
		var g errgroup.Group
		amps := make([]*amp, 0, ampCount)
		for i := 0; i < ampCount; i++ {
			amp := newAmp(i, p[i])
			if i == 0 {
				amp.in <- 0 // first amp phase
			} else {
				amp.in = amps[i-1].out
			}
			amps = append(amps, amp)
			g.Go(func() error { return amp.run(program) })
		}

		g.Wait()

		output := <-amps[len(amps)-1].out
		if output > maxOutputSignal {
			maxOutputSignal = output
			copy(maxPhaseSettings, p)
		}
	}, 0)

	fmt.Println(maxPhaseSettings)
	fmt.Println(maxOutputSignal)
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
		case 1:
			return <-a.in, nil
		default:
			return 0, fmt.Errorf("unknown input state %d", inputState)
		}
	}

	output := func(x int) error {
		a.out <- x
		return nil
	}

	p := make([]int, len(program))
	copy(p, program)
	return intcode.Run(p, input, output)
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
