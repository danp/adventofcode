package main

import (
	"fmt"

	"github.com/danp/adventofcode/scaffold"
)

func main() {
	lines := scaffold.Lines()

	var valid, validOne int
	for _, l := range lines {
		var pwp passwordPolicy
		if _, err := fmt.Sscanf(l, "%d-%d %c: %s", &pwp.min, &pwp.max, &pwp.letter, &pwp.password); err != nil {
			panic(err)
		}

		fmt.Println(pwp)
		if pwp.isValid() {
			valid++
		}
		if pwp.isValidOneLetter() {
			validOne++
		}

	}
	fmt.Println(valid, validOne)
}

type passwordPolicy struct {
	letter   rune
	min, max int
	password string
}

func (p passwordPolicy) isValid() bool {
	var count int
	for _, c := range p.password {
		if c == p.letter {
			count++
		}
		if count > p.max {
			return false
		}
	}
	return count >= p.min
}

func (p passwordPolicy) isValidOneLetter() bool {
	return (rune(p.password[p.min-1]) == p.letter) != (rune(p.password[p.max-1]) == p.letter)
}

func (p passwordPolicy) String() string {
	return fmt.Sprintf("%d-%d %c: %q valid=%t validOne=%t", p.min, p.max, p.letter, p.password, p.isValid(), p.isValidOneLetter())
}
