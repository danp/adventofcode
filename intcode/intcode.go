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
	x    func(*vm) error
}

var halt = errors.New("halt")

var ops = map[int]op{
	1: {
		name: "add",
		code: 1,
		pc:   3,
		x: func(v *vm) error {
			v.set(2, v.mval(0)+v.mval(1))
			return nil
		},
	},
	2: {
		name: "mult",
		code: 2,
		pc:   3,
		x: func(v *vm) error {
			v.set(2, v.mval(0)*v.mval(1))
			return nil
		},
	},
	3: {
		name: "input",
		code: 3,
		pc:   1,
		x: func(v *vm) error {
			if v.input == nil {
				return errors.New("program wants input but no input func provided")
			}
			in, err := v.input()
			if err != nil {
				return err
			}
			v.set(0, in)
			return nil
		},
	},
	4: {
		name: "output",
		code: 4,
		pc:   1,
		x: func(v *vm) error {
			if v.output == nil {
				return errors.New("program wants to output but no output func provided")
			}
			return v.output(v.mval(0))
		},
	},
	5: {
		name: "jump-if-true",
		code: 5,
		pc:   2,
		x: func(v *vm) error {
			if v.mval(0) > 0 {
				v.jump(v.mval(1))
			}
			return nil
		},
	},
	6: {
		name: "jump-if-false",
		code: 6,
		pc:   2,
		x: func(v *vm) error {
			if v.mval(0) == 0 {
				v.jump(v.mval(1))
			}
			return nil
		},
	},
	7: {
		name: "less-than",
		code: 7,
		pc:   3,
		x: func(v *vm) error {
			var res int
			if v.mval(0) < v.mval(1) {
				res = 1
			}
			v.set(2, res)
			return nil
		},
	},
	8: {
		name: "equals",
		code: 8,
		pc:   3,
		x: func(v *vm) error {
			var res int
			if v.mval(0) == v.mval(1) {
				res = 1
			}
			v.set(2, res)
			return nil
		},
	},
	9: {
		name: "adjust-relative-base",
		code: 9,
		pc:   1,
		x: func(v *vm) error {
			v.relbase += v.mval(0)
			return nil
		},
	},
	99: {
		name: "halt",
		code: 99,
		pc:   0,
		x: func(v *vm) error {
			return halt
		},
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
	relative  pmode = 3
)

type vm struct {
	program []int
	mem     []int
	pos     int
	ins     instruction

	input  func() (int, error)
	output func(int) error

	jumped  bool
	relbase int
}

func (v *vm) run() error {
	copy(v.mem, v.program)

	for v.pos <= len(v.mem) {
		if err := v.stepInstruction(); err != nil {
			return err
		}

		if err := v.ins.op.x(v); err != nil {
			if err == halt {
				err = nil
			}
			return err
		}
	}

	panic("got here")
}

func (v *vm) stepInstruction() error {
	if v.ins.op.code > 0 && !v.jumped {
		v.pos += v.ins.op.pc
	}
	v.jumped = false

	ins, err := parseInstruction(v.val(0))
	if err != nil {
		return err
	}

	v.ins = ins
	v.pos++

	return nil
}

func (v *vm) val(i int) int {
	return v.mem[v.pos+i]
}

func (v *vm) mval(i int) int {
	switch v.ins.pmodes[i] {
	case position:
		return v.mem[v.val(i)]
	case immediate:
		return v.val(i)
	case relative:
		return v.mem[v.relbase+v.val(i)]
	default:
		panic("unknown pmode")
	}
}

func (v *vm) set(i, val int) {
	j := v.val(i)
	if v.ins.pmodes[i] == relative {
		j += v.relbase
	}
	v.mem[j] = val
}

func (v *vm) jump(pos int) {
	v.pos = pos
	v.jumped = true
}

// Run runs program, working in mem, calling input when input
// is requested and calling output when output is requested.
//
// The contents of mem will be modified.
func Run(program, mem []int, input func() (int, error), output func(int) error) error {
	vm := &vm{program: program, mem: mem, input: input, output: output}
	return vm.run()
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
		case 2:
			m = relative
		default:
			return ins, fmt.Errorf("unknown param mode in instruction %d", in)
		}

		ins.pmodes = append(ins.pmodes, m)

		in /= 10
	}

	return ins, nil
}
