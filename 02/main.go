package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func main() {
	switch os.Args[1] {
	case "one":
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}

		program, err := parse(string(b))
		if err != nil {
			panic(err)
		}

		if err := run(program); err != nil {
			panic(err)
		}

		fmt.Println(program[0])
	case "finder":
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}

		program, err := parse(string(b))
		if err != nil {
			panic(err)
		}

		for {
			p := make([]int, len(program))
			copy(p, program)

			noun := rand.Intn(100)
			verb := rand.Intn(100)

			p[1] = noun
			p[2] = verb

			if err := run(p); err != nil {
				panic(err)
			}

			if p[0] == 19690720 {
				fmt.Println(noun, verb, 100*noun+verb)
				break
			}
		}
	}
}

func run(program []int) error {
	var pos int
	for pos <= len(program) {
		op := program[pos]

		switch op {
		case 1:
			aPos := program[pos+1]
			bPos := program[pos+2]
			oPos := program[pos+3]
			a := program[aPos]
			b := program[bPos]
			program[oPos] = a + b
			pos += 4
		case 2:
			aPos := program[pos+1]
			bPos := program[pos+2]
			oPos := program[pos+3]
			a := program[aPos]
			b := program[bPos]
			program[oPos] = a * b
			pos += 4
		case 99:
			return nil
		default:
			return fmt.Errorf("unknown opcode %d at pos %d", op, pos)
		}
	}

	panic("got here")
}

func parse(input string) ([]int, error) {
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
