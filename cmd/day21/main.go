package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

	if os.Args[1] == "part1" {
		part1(program)
	} else {
		part2(program)
	}
}

func part1(program []int) {
	scripts := []string{
		// NOT D J // ? ? ? !d -> useless
		"NOT A J",                         // !a ? ? ?
		"NOT B J",                         // ? !b ? ?
		"NOT C J",                         // ? ? !c ?
		"OR A T\nOR B J\nOR T J\nAND D J", // !a or !b ? d
		"OR B T\nOR C J\nOR T J\nAND D J", // ? !b or !c d
		"OR A T\nAND B T\nAND D T\nNOT C J\nAND T J",           // a b !c d
		"OR A T\nAND C T\nAND D T\nNOT B J\nAND T J",           // a !b c d
		"OR B T\nAND C T\nAND D T\nNOT A J\nAND T J",           // !a b c d
		"NOT B T\nNOT C J\nAND T J\nAND A J\nAND D J",          // a !b !c d
		"NOT A T\nNOT C J\nAND T J\nAND B J\nAND D J",          // !a b !c d
		"NOT A T\nNOT B J\nAND T J\nNOT C T\nAND T J\nAND D J", // !a !b !c d
		"OR A J\nOR B J\nOR C J\nAND D J",                      // !a or !b or !c and d
		"NOT A J\nNOT C T\nOR T J",                             // !a or !c ? ?
		"NOT A J\nNOT C T\nOR T J\nOR B T\nOR D T\nAND T J",    // !a or !c and b or d
		"NOT A J\nNOT B T\nOR T J\nNOT C T\nOR T J\nAND D J",   // (!a or !b or !c) and d
	}

	var done bool
	for _, script := range scripts {
		var inb bytes.Buffer
		fmt.Fprintln(&inb, strings.TrimSpace(script))
		fmt.Fprintln(&inb, "WALK")
		input := func() (int, error) {
			c, err := inb.ReadByte()
			return int(c), err
		}

		var outb bytes.Buffer
		output := func(x int) error {
			if x > 255 {
				fmt.Println(x)
				done = true
			} else {
				fmt.Fprint(&outb, x)
			}
			return nil
		}

		mem := make([]int, 3072)
		if err := intcode.Run(program, mem, input, output); err != nil {
			panic(err)
		}

		if done {
			fmt.Println()
			fmt.Println(outb.String())
			break
		}

		showlines(&outb)
	}
}

func part2(program []int) {
	scripts := []string{
		"NOT A J", // !a ? ? ?
		"NOT A J\nNOT B T\nOR T J\nNOT C T\nOR T J\nAND D J", // (!a or !b or !c) and d
		"OR D J\nOR E J\nAND H J",
		"OR D J\nOR E J\nAND H J\nAND C J",
		"OR E J\nAND D J\nAND F J",
		"NOT B J\nNOT C T\nOR T J\nAND H J\nNOT A T\nOR T J\nAND D J",
	}

	var done bool
	for _, script := range scripts {
		var inb bytes.Buffer
		fmt.Fprintln(&inb, strings.TrimSpace(script))
		fmt.Fprintln(&inb, "RUN")
		input := func() (int, error) {
			c, err := inb.ReadByte()
			return int(c), err
		}

		var outb bytes.Buffer
		output := func(x int) error {
			if x > 255 {
				fmt.Println(x)
				done = true
			} else {
				fmt.Fprint(&outb, string(rune(x)))
			}
			return nil
		}

		mem := make([]int, 3072)
		if err := intcode.Run(program, mem, input, output); err != nil {
			panic(err)
		}

		if done {
			fmt.Println()
			fmt.Println(outb.String())
			break
		}

		showlines(&outb)
	}
}

func showlines(b *bytes.Buffer) {
	outpics := strings.Split(strings.TrimSpace(b.String()), "\n\n")

	lpiclines := strings.Split(strings.TrimSpace(outpics[len(outpics)-1]), "\n")
	lpicline := lpiclines[len(lpiclines)-1]
	var stidx int

	for i, op := range outpics[3:] {
		piclines := strings.Split(strings.TrimSpace(op), "\n")
		if strings.Index(piclines[1], "@") != -1 {
			x := strings.Split(outpics[3+i-1], "\n")
			stidx = strings.Index(x[2], "@")
			fmt.Println(x[2])
			break
		}
	}

	fmt.Println(lpicline)
	fmt.Println(strings.Repeat(" ", stidx+1) + "ABCDEFGHI")
}
