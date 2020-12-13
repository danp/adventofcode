package scaffold

import "strconv"

func Ints(lines []string) []int {
	ints := make([]int, 0, len(lines))
	for _, l := range lines {
		ints = append(ints, Int(l))
	}
	return ints
}

func Int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
