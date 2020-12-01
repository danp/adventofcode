package scaffold

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func Lines() []string {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	return strings.Split(strings.TrimSpace(string(b)), "\n")
}

func Ints(lines []string) []int {
	ints := make([]int, 0, len(lines))
	for _, l := range lines {
		i, err := strconv.Atoi(l)
		if err != nil {
			panic(err)
		}
		ints = append(ints, i)
	}
	return ints
}
