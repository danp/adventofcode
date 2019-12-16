package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	digits := convert(strings.TrimSpace(string(b)))

	if os.Args[1] == "part1" {
		for i := 0; i < 100; i++ {
			phase(digits)
		}

		fmt.Println(digits[:8])
		// 25131128
		return
	}

	offset, err := strconv.Atoi(string(b)[:7])
	if err != nil {
		panic(err)
	}
	fmt.Println(offset)

	for i := 1; i < 10000; i++ {
		digits = append(digits, convert(strings.TrimSpace(string(b)))...)
	}

	digits = digits[offset:]

	fmt.Println(len(digits))

	for i := 0; i < 100; i++ {
		sum := 0
		for j := len(digits) - 1; j >= 0; j-- {
			sum += digits[j]
			digits[j] = sum % 10
		}
	}

	fmt.Println(join(digits[:8]))

}

func phase(data []int) {
	for i := 0; i < len(data); i++ {
		var accum int

		for j := i; j < len(data); j++ {
			m := multiplier(i, j)
			accum += data[j] * m
		}
		if accum <= -0 {
			accum *= -1
		}
		if accum > 9 {
			accum %= 10
		}

		data[i] = accum
	}
}

func multiplier(opos, ipos int) int {
	ipos++ // shift left

	repeats := opos + 1                       // 0,0,0,0,...
	cycle := repeats * 4                      // 0,1,-1,0,...
	cnum := ipos / cycle                      // what cycle we are in
	cpos := (ipos - (cycle * cnum)) / repeats // what cycle position we are in

	// fmt.Println("repeating", repeats, "times, cycle length is", cycle, "and", ipos, "is in cycle", cnum, "position", cpos)

	return []int{0, 1, 0, -1}[cpos]
}

func convert(s string) []int {
	out := make([]int, len(s))
	for i, c := range s {
		x, err := strconv.Atoi(string(c))
		if err != nil {
			panic(err)
		}
		out[i] = x
	}
	return out
}

func join(is []int) string {
	var b strings.Builder
	for _, i := range is {
		b.WriteString(strconv.Itoa(i))
	}
	return b.String()
}
