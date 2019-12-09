package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/danp/adventofcode2019/intcode"
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

	mem := make([]int, len(program))

	switch os.Args[1] {
	case "ac":
		input := func() (int, error) { return 1, nil }
		output := func(x int) error { fmt.Println(x); return nil }
		if err := intcode.Run(program, mem, input, output); err != nil {
			panic(err)
		}
	case "radiators":
		input := func() (int, error) { return 5, nil }
		output := func(x int) error { fmt.Println(x); return nil }
		if err := intcode.Run(program, mem, input, output); err != nil {
			panic(err)
		}
	}
}
