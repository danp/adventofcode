package intcode

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type op struct {
	name string
	code int
	pc   int
}

var ops = map[int]op{
	1: {
		name: "add",
		code: 1,
		pc:   3,
	},
	2: {
		name: "mult",
		code: 2,
		pc:   3,
	},
	3: {
		name: "input",
		code: 3,
		pc:   1,
	},
	4: {
		name: "output",
		code: 4,
		pc:   1,
	},
	5: {
		name: "jump-if-true",
		code: 5,
		pc:   2,
	},
	6: {
		name: "jump-if-false",
		code: 6,
		pc:   2,
	},
	7: {
		name: "less-than",
		code: 7,
		pc:   3,
	},
	8: {
		name: "equals",
		code: 8,
		pc:   3,
	},
	99: {
		name: "halt",
		code: 99,
		pc:   0,
	},
}

type instruction struct {
	op     op
	pmodes []pmode
}

type pmode int

const (
	position  pmode = 1
	immediate pmode = 2
)

// Run runs program, calling input when input is requested and
// calling output when output is requested.
// The contents of program may be modified.
func Run(program []int, input func() (int, error), output func(int) error) error {
	var pos int

	for pos <= len(program) {
		inst, err := parseInstruction(program[pos])
		if err != nil {
			return err
		}

		val := func(i int) int {
			return program[pos+1+i]
		}

		mval := func(i int) int {
			v := val(i)
			if inst.pmodes[i] == position {
				v = program[v]
			}
			return v
		}

		switch inst.op.code {
		case 1:
			program[val(2)] = mval(0) + mval(1)
		case 2:
			program[val(2)] = mval(0) * mval(1)
		case 3:
			if input == nil {
				return errors.New("program wants input but no input func provided")
			}
			in, err := input()
			if err != nil {
				return err
			}
			program[val(0)] = in
		case 4:
			if output == nil {
				return errors.New("program wants to output but no output func provided")
			}
			if err := output(mval(0)); err != nil {
				return err
			}
		case 5:
			if mval(0) > 0 {
				pos = mval(1)
				continue
			}
		case 6:
			if mval(0) == 0 {
				pos = mval(1)
				continue
			}
		case 7:
			var res int
			if mval(0) < mval(1) {
				res = 1
			}
			program[val(2)] = res
		case 8:
			var res int
			if mval(0) == mval(1) {
				res = 1
			}
			program[val(2)] = res
		case 99:
			return nil
		}

		pos += 1 + inst.op.pc
	}

	panic("got here")
}

// Parse takes a program string in the form `1,2,3,...` and returns a
// slice of int ready for use with Run.
func Parse(input string) ([]int, error) {
	input = strings.ReplaceAll(input, "\n", "")
	parts := strings.Split(input, ",")

	var program []int
	for _, p := range parts {
		i, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		program = append(program, i)
	}

	return program, nil
}

func buildOut(program []int) string {
	var out string

	for i, p := range program {
		if i > 0 {
			out += ","
		}
		out += strconv.Itoa(p)
	}

	return out
}

func parseInstruction(in int) (instruction, error) {
	var ins instruction

	opcode := in % 100
	op, ok := ops[opcode]
	if !ok {
		return ins, fmt.Errorf("unknown opcode %d in instruction %d", opcode, in)
	}

	ins.op = op
	in /= 100

	for i := 0; i < op.pc; i++ {
		var m pmode
		switch in % 10 {
		case 0:
			m = position
		case 1:
			m = immediate
		default:
			return ins, fmt.Errorf("unknown param mode in instruction %d", in)
		}

		ins.pmodes = append(ins.pmodes, m)

		in /= 10
	}

	return ins, nil
}
