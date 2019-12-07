package main

import (
	"fmt"
	"os"
)

func main() {
	var valids []int

	check := checkPassword
	if len(os.Args) > 1 && os.Args[1] == "harder" {
		check = checkPasswordHarder
	}

	for i := 271973; i <= 785961; i++ {
		if check(i) {
			valids = append(valids, i)
		}
	}

	fmt.Println(len(valids))
}

func checkPassword(pw int) bool {
	last := -1
	var sawDouble bool

	for pw > 0 {
		d := pw % 10
		pw /= 10

		if last == -1 {
			last = d
			continue
		}

		if d > last {
			return false
		}

		if d == last {
			sawDouble = true
		}

		last = d
	}

	return sawDouble
}

func checkPasswordHarder(pw int) bool {
	last := -1
	runCount := 1
	var sawDouble bool

	for pw > 0 {
		d := pw % 10
		pw /= 10

		if last == -1 {
			last = d
			continue
		}

		if d > last {
			return false
		}

		if d == last {
			runCount++
		} else {
			if runCount == 2 {
				sawDouble = true
			}
			runCount = 1
		}

		last = d
	}

	return sawDouble || runCount == 2
}
