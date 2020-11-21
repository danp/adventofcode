package main

import (
	"fmt"
	"io/ioutil"
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

	mem := make([]int, 2048)
	var in int

	switch os.Args[1] {
	case "boost":
		in = 1
	case "distress":
		in = 2
	}

	input := func() (int, error) { return in, nil }
	output := func(x int) error { fmt.Println(x); return nil }

	if err := intcode.Run(program, mem, input, output); err != nil {
		panic(err)
	}
}
