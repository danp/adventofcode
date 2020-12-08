package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	prog, err := parse(lines)
	if err != nil {
		panic(err)
	}

	seenpc := make(map[int]bool)
	m1 := machine{prog: prog}
	for !seenpc[m1.pc] {
		seenpc[m1.pc] = true
		m1.next()
	}

	fmt.Println("acc before executing an instruction a second time:", m1.acc)

nextHist:
	for i := len(m1.hist) - 1; i >= 0; i-- {
		h := m1.hist[i]
		if h.op != "nop" && h.op != "jmp" {
			continue
		}

		prog2 := make([]inst, len(prog))
		copy(prog2, prog)
		if h.op == "nop" {
			prog2[h.pc].op = "jmp"
		} else {
			prog2[h.pc].op = "nop"
		}

		seenpc = make(map[int]bool)
		m2 := machine{prog: prog2}
		for {
			if seenpc[m2.pc] {
				continue nextHist
			}
			seenpc[m2.pc] = true

			if !m2.next() {
				fmt.Println("acc after successful termination of fixed program:", m2.acc)
				break nextHist
			}
		}
	}
}

type inst struct {
	pc  int
	op  string
	arg int
}

func parse(prog []string) ([]inst, error) {
	out := make([]inst, 0, len(prog))

	for i, p := range prog {
		f := strings.Fields(p)
		op := f[0]
		arg64, err := strconv.ParseInt(f[1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("line %d: %s", i+1, err)
		}
		arg := int(arg64)

		out = append(out, inst{pc: i, op: op, arg: arg})
	}

	return out, nil
}

type machine struct {
	prog []inst
	pc   int
	acc  int

	hist []inst
}

func (m *machine) next() bool {
	ins := m.prog[m.pc]

	m.hist = append(m.hist, ins)

	switch ins.op {
	case "nop":
		m.pc++
	case "acc":
		m.acc += ins.arg
		m.pc++
	case "jmp":
		m.pc += ins.arg
	}

	if m.pc < 0 {
		panic("pc < 0")
	}

	return m.pc < len(m.prog)
}
