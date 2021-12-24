package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	var cmds []cmd
	lines := scaffold.Lines()
	for _, l := range lines {
		c := parseCmd(l)
		cmds = append(cmds, c)
	}

	try := func(n int) {
		buf := fmt.Sprint(n)
		inp := func() int {
			n := scaffold.Int(buf[0:1])
			buf = buf[1:]
			fmt.Printf("n: %v\n", n)
			return n
		}
		a := newALU(inp)
		for _, c := range cmds {
			a.run(c)
		}
	}

	for i := 0; i < 10; i++ {
		// 11111111111111 - 99999999999999, no zeros
		n := rand.Int63n(93579239999998)
	try()
}

type alu struct {
	vars map[string]int

	inp func() int
}

func newALU(inp func() int) *alu {
	return &alu{vars: make(map[string]int), inp: inp}
}

func (a *alu) run(c cmd) {
	value := func(arg string) int {
		switch arg {
		case "w", "x", "y", "z":
			return a.vars[arg]
		default:
			return scaffold.Int(arg)
		}
	}

	switch c.op {
	case "inp":
		fmt.Printf("a.vars: %v\n", a.vars)
		a.vars[c.args[0]] = a.inp()
	case "add":
		a.vars[c.args[0]] = a.vars[c.args[0]] + value(c.args[1])
	case "mul":
		a.vars[c.args[0]] = a.vars[c.args[0]] * value(c.args[1])
	case "div":
		a.vars[c.args[0]] = a.vars[c.args[0]] / value(c.args[1])
	case "mod":
		a.vars[c.args[0]] = a.vars[c.args[0]] % value(c.args[1])
	case "eql":
		var v int
		if a.vars[c.args[0]] == value(c.args[1]) {
			v = 1
		}
		a.vars[c.args[0]] = v
	}
}

type cmd struct {
	op   string
	args []string
}

func parseCmd(s string) cmd {
	f := strings.Fields(s)
	return cmd{f[0], f[1:]}
}
