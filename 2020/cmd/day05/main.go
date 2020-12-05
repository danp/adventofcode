package main

import (
	"fmt"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	for _, l := range lines {
		row := ranger{min: 0, max: 127}
		for _, c := range l[:7] {
			switch c {
			case 'F':
				row = row.left()
			case 'B':
				row = row.right()
			}
		}
		seat := ranger{min: 0, max: 7}
		for _, c := range l[7:] {
			switch c {
			case 'L':
				seat = seat.left()
			case 'R':
				seat = seat.right()
			}
		}

		id := row.max*8 + seat.max
		fmt.Println(l, "row", row.max, "seat", seat.max, "id", id)
	}

	// part 1 solved with:
	// go run main.go < input| awk '{ print $NF }' | sort -rn | head -n1
	// yielding 896

	// part 2 solved with:
	// diff -u <(seq 53 896) <(go run main.go < input| awk '{ print $NF }' | sort -n)
	// yielding a single removed line for id 659
}

type ranger struct {
	min, max int
}

func (r ranger) half() int {
	return r.min + ((r.max - r.min) / 2)
}

func (r ranger) left() ranger {
	return ranger{min: r.min, max: r.half()}
}

func (r ranger) right() ranger {
	return ranger{min: r.half(), max: r.max}
}
