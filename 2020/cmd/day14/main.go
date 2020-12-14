package main

import (
	"fmt"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	part1(lines)
	part2(lines)
}

func part1(lines []string) {
	// address -> val
	mem := make(map[int]int)

	// bit -> 1/0
	mask := make(map[int]int)

	for _, l := range lines {
		f := strings.Fields(l)
		if f[0] == "mask" {
			for k := range mask {
				delete(mask, k)
			}
			for i := len(f[2]) - 1; i >= 0; i-- {
				c := f[2][i]
				if c == 'X' {
					continue
				}

				bit := len(f[2]) - i - 1
				var val int
				if c == '1' {
					val = 1
				}
				mask[bit] = val
			}
			fmt.Println(mask)
			continue
		}

		var addr, val int
		if _, err := fmt.Sscanf(l, "mem[%d] = %d", &addr, &val); err != nil {
			panic(err)
		}

		mval := val
		for b, v := range mask {
			bit := 1 << b
			if v == 1 {
				mval |= bit
			} else if mval&bit != 0 {
				mval ^= bit
			}
		}
		fmt.Println("write", addr, "=", mval)
		mem[addr] = mval
	}

	fmt.Println(mem)

	sum := 0
	for _, v := range mem {
		sum += v
	}
	fmt.Println(sum)
}

func part2(lines []string) {
	// address -> val
	mem := make(map[int]int)

	// bit -> 1 to set to 1 or 2 for floating
	mask := make(map[int]int)

	for _, l := range lines {
		f := strings.Fields(l)
		if f[0] == "mask" {
			for k := range mask {
				delete(mask, k)
			}
			for i := len(f[2]) - 1; i >= 0; i-- {
				c := f[2][i]
				if c == '0' {
					continue
				}
				bit := len(f[2]) - i - 1

				var val int
				switch c {
				case '1':
					val = 1
				case 'X':
					val = 2
				}
				mask[bit] = val
			}
			fmt.Println(mask)
			continue
		}

		var addr, val int
		if _, err := fmt.Sscanf(l, "mem[%d] = %d", &addr, &val); err != nil {
			panic(err)
		}

		fmt.Println(mask, addr)
		addrs := addrs(mask, addr)
		fmt.Println(addrs)

		for _, a := range addrs {
			fmt.Println("write", a, "=", val)
			mem[a] = val
		}
	}

	fmt.Println(mem)

	sum := 0
	for _, v := range mem {
		sum += v
	}
	fmt.Println(sum)
}

func addrs(mask map[int]int, addr int) []int {
	var out []int
	var floats []int

	for b, v := range mask {
		switch v {
		case 1:
			addr |= 1 << b
		case 2:
			floats = append(floats, b)
		}
	}

	if len(floats) == 0 {
		return []int{addr}
	}

	for _, p := range bitperms(floats) {
		fmt.Println("perm", p)
		addr := addr
		for b, v := range p {
			bit := 1 << b
			switch v {
			case 0:
				if addr&bit != 0 {
					addr ^= bit
				}
			case 1:
				addr |= 1 << b
			}
		}
		out = append(out, addr)
	}

	return out
}

func bitperms(bits []int) []map[int]int {
	var out []map[int]int

	for i := 0; i < 1<<len(bits); i++ {
		m := make(map[int]int)
		for j, b := range bits {
			v := 0
			if i&(1<<j) != 0 {
				v = 1
			}
			m[b] = v
		}
		out = append(out, m)
	}

	return out
}
