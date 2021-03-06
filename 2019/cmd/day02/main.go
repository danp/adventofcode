package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/danp/adventofcode/2019/intcode"
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

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "need an arg of one or finder")
		os.Exit(1)
	}

	mem := make([]int, len(program))

	switch os.Args[1] {
	case "one":
		if err := intcode.Run(program, mem, nil, nil); err != nil {
			panic(err)
		}

		fmt.Println(mem[0])
	case "finder":
		for {
			p := make([]int, len(program))
			copy(p, program)

			noun := rand.Intn(100)
			verb := rand.Intn(100)

			p[1] = noun
			p[2] = verb

			if err := intcode.Run(p, mem, nil, nil); err != nil {
				panic(err)
			}

			if mem[0] == 19690720 {
				fmt.Println(noun, verb, 100*noun+verb)
				break
			}
		}
	}
}
