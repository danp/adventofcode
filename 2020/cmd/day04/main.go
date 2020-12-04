package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	strict := false
	if len(os.Args) > 1 && os.Args[1] == "part2" {
		strict = true
	}

	fmt.Println(check(lines, strict))
}

func check(lines []string, strict bool) int {
	var valid int
	data := make(map[string]string)
	for _, l := range lines {
		if l == "" {
			if isValid(data, strict) {
				valid++
				fmt.Println("valid  ", data)
			} else {
				fmt.Println("invalid", data)
			}
			data = make(map[string]string)
		}
		for _, f := range strings.Fields(l) {
			parts := strings.SplitN(f, ":", 2)
			data[parts[0]] = parts[1]
		}
	}

	if isValid(data, strict) {
		fmt.Println("valid  ", data)
		valid++
	} else {
		fmt.Println("invalid", data)
	}
	return valid
}

func checkIntRange(s string, min, max int) bool {
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return false
	}
	ni := int(n)
	return ni >= min && ni <= max
}

func inSet(s string, set []string) bool {
	for _, ss := range set {
		if s == ss {
			return true
		}
	}
	return false
}

var (
	hclMatch = regexp.MustCompile(`\A#[0-9a-f]{6}\z`)
	pidMatch = regexp.MustCompile(`\A\d{9}\z`)
)

func isValid(data map[string]string, strict bool) bool {
	if !strict && len(data) == 8 {
		return true
	}

	if len(data) < 7 {
		return false
	}

	if !strict {
		_, ok := data["cid"]
		return !ok
	}

	if !checkIntRange(data["byr"], 1920, 2002) {
		return false
	}

	if !checkIntRange(data["iyr"], 2010, 2020) {
		return false
	}

	if !checkIntRange(data["eyr"], 2020, 2030) {
		return false
	}

	hgt := data["hgt"]
	if strings.HasSuffix(hgt, "cm") {
		if !checkIntRange(hgt[:len(hgt)-2], 150, 193) {
			return false
		}
	} else if strings.HasSuffix(hgt, "in") {
		if !checkIntRange(hgt[:len(hgt)-2], 59, 76) {
			return false
		}
	} else {
		return false
	}

	if !inSet(data["ecl"], []string{"amb", "blu", "brn", "gry", "grn", "hzl", "oth"}) {
		return false
	}

	if !hclMatch.MatchString(data["hcl"]) {
		return false
	}

	if !pidMatch.MatchString(data["pid"]) {
		return false
	}

	return true
}
