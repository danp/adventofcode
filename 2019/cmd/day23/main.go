package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

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

	cs := make([]*computer, 50)

	var (
		nMu sync.Mutex
		ls  = time.Now()
		nm  message
		ny  = make(map[int]int)
	)

	ns := func(m message) {
		nMu.Lock()
		defer nMu.Unlock()
		ls = time.Now()

		if m.tn == 255 {
			nm = m
		} else {
			cs[m.tn].recv(m)
		}
	}

	go func() {
		tick := time.NewTicker(time.Second)
		defer tick.Stop()

		for range tick.C {
			nMu.Lock()
			since := time.Since(ls)
			nMu.Unlock()
			if since > time.Second {
				fmt.Println("nat sending", nm)
				nMu.Lock()
				ny[nm.y]++
				if ny[nm.y] > 1 {
					fmt.Println("nat dup y", nm.y)
				}
				nm.tn = 0
				nMu.Unlock()
				ns(nm)
			}
		}
	}()

	var g errgroup.Group

	for i := 0; i < len(cs); i++ {
		cs[i] = &computer{n: i, networkSend: ns}
	}

	for i := 0; i < len(cs); i++ {
		i := i
		g.Go(func() error { return cs[i].run(program) })
	}

	if err := g.Wait(); err != nil {
		panic(err)
	}
}

type message struct {
	tn   int
	x, y int
}

type computer struct {
	n int

	networkSend func(m message)

	rMu sync.Mutex
	rq  []message
}

func (c *computer) recv(m message) {
	c.rMu.Lock()
	defer c.rMu.Unlock()
	c.rq = append(c.rq, m)
}

func (c *computer) run(program []int) error {
	var (
		ist  int
		rm   message
		rmst int
	)
	input := func() (int, error) {
		if ist == 0 {
			ist++
			return c.n, nil
		}

		c.rMu.Lock()
		defer c.rMu.Unlock()

		switch rmst {
		case 0:
			if len(c.rq) == 0 {
				return -1, nil
			}
			rm = c.rq[0]
			c.rq = c.rq[1:]
			rmst++
			return rm.x, nil
		case 1:
			rmst = 0
			return rm.y, nil
		}

		panic("got here")
	}

	var (
		ost int
		om  message
	)
	output := func(x int) error {
		switch ost {
		case 0:
			om.tn = x
			ost++
		case 1:
			om.x = x
			ost++
		case 2:
			om.y = x
			c.networkSend(om)
			ost = 0
		}

		return nil
	}

	mem := make([]int, 4096)
	return intcode.Run(program, mem, input, output)
}
