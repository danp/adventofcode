package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	if os.Args[1] == "part1" {
		cs, err := parseCommands(strings.TrimSpace(string(b)))
		if err != nil {
			panic(err)
		}

		deck := newDeck(10007)
		for _, c := range cs {
			deck = c(deck)
		}

		for i := 0; i < len(deck); i++ {
			if deck[i] == 2019 {
				fmt.Println(i)
			}
		}
	} else {
		cs, err := parseCommands2(strings.TrimSpace(string(b)))
		if err != nil {
			panic(err)
		}

		st := newState(119315717514047, 101741582076661)
		for _, c := range cs {
			c.f(st)
		}
		fmt.Println(st.get(2020))
	}
}

func newDeck(n int) []int {
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = i
	}
	return out
}

type state struct {
	c, r   *big.Int
	incr   *big.Int
	offset *big.Int
}

func newState(cards, repeats int) *state {
	return &state{c: bi(cards), r: bi(repeats), incr: bi(1), offset: bi(0)}
}

func b0() *big.Int {
	return bi(0)
}

func bi(n int) *big.Int {
	return big.NewInt(int64(n))
}

func (s *state) get(pos int) *big.Int {
	incr := b0().Exp(s.incr, s.r, s.c)
	offset := b0().Sub(bi(1), incr)
	invmod := b0().Exp(b0().Sub(bi(1), s.incr), b0().Sub(s.c, bi(2)), s.c)
	offset.Mul(offset, invmod)
	offset.Mul(offset, s.offset)
	x := b0().Mul(bi(pos), incr)
	x.Add(x, offset)
	x.Mod(x, s.c)
	return x
}

func (s *state) deal() {
	s.incr.Mul(s.incr, bi(-1))
	s.offset.Add(s.offset, s.incr)
}

func (s *state) cut(n int) {
	s.offset.Add(s.offset, b0().Mul(bi(n), s.incr))
}

func (s *state) dealIncrement(n int) {
	s.incr.Mul(s.incr, b0().Exp(bi(n), b0().Sub(s.c, bi(2)), s.c))
}

func dealNewStack(stack []int) []int {
	out := make([]int, len(stack))
	for i := len(stack) - 1; i >= 0; i-- {
		out[i] = stack[len(stack)-i-1]
	}
	return out
}

func cutN(deck []int, n int) []int {
	out := make([]int, len(deck))
	if n < 0 {
		n *= -1
		for i := 0; i < n; i++ {
			out[i] = deck[len(deck)-n+i]
		}
		for i := n; i < len(deck); i++ {
			out[i] = deck[i-n]
		}
	} else {
		for i := n; i < len(deck); i++ {
			out[i-n] = deck[i]
		}
		for i := 0; i < n; i++ {
			out[len(deck)-n+i] = deck[i]
		}
	}
	return out
}

func dealIncrement(deck []int, n int) []int {
	out := make([]int, len(deck))

	for i := 0; i < len(deck); i++ {
		opos := i * n
		round := opos / len(deck)
		opos -= round * len(deck)
		out[opos] = deck[i]
	}

	return out
}

var (
	cutRe           = regexp.MustCompile(`\Acut (-?\d+)\z`)
	dealIncrementRe = regexp.MustCompile(`\Adeal with increment (\d+)\z`)
)

func parseCommands(input string) ([]func([]int) []int, error) {
	lines := strings.Split(input, "\n")
	out := make([]func([]int) []int, len(lines))
	for i, l := range lines {
		if l == "deal into new stack" {
			out[i] = dealNewStack
		} else if m := cutRe.FindStringSubmatch(l); len(m) > 1 {
			n, err := strconv.Atoi(m[1])
			if err != nil {
				return nil, err
			}
			out[i] = func(d []int) []int { return cutN(d, n) }
		} else if m := dealIncrementRe.FindStringSubmatch(l); len(m) > 1 {
			n, err := strconv.Atoi(m[1])
			if err != nil {
				return nil, err
			}
			out[i] = func(d []int) []int { return dealIncrement(d, n) }
		} else {
			return nil, fmt.Errorf("unknown command %q", l)
		}
	}
	return out, nil
}

type command struct {
	s string
	f func(*state)
}

func parseCommands2(input string) ([]command, error) {
	lines := strings.Split(input, "\n")
	out := make([]command, len(lines))

	for i, l := range lines {
		c := command{s: l}

		if l == "deal into new stack" {
			c.f = func(st *state) {
				st.deal()
			}
		} else if m := cutRe.FindStringSubmatch(l); len(m) > 1 {
			n, err := strconv.Atoi(m[1])
			if err != nil {
				return nil, err
			}
			c.f = func(st *state) {
				st.cut(n)
			}
		} else if m := dealIncrementRe.FindStringSubmatch(l); len(m) > 1 {
			n, err := strconv.Atoi(m[1])
			if err != nil {
				return nil, err
			}
			c.f = func(st *state) {
				st.dealIncrement(n)
			}
		} else {
			return nil, fmt.Errorf("unknown command %q", l)
		}

		out[i] = c
	}
	return out, nil
}

func s(stack []int) string {
	return fmt.Sprintf("%v", stack)
}
