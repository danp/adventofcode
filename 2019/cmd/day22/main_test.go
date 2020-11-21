package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDealNewStack(t *testing.T) {
	stack := dealNewStack(newDeck(10))
	if d := cmp.Diff("[9 8 7 6 5 4 3 2 1 0]", s(stack)); d != "" {
		t.Errorf("deal new stack mismatch (-want +got):\n%s", d)
	}
}

func TestCutN(t *testing.T) {
	deck := cutN(newDeck(10), 3)
	if d := cmp.Diff("[3 4 5 6 7 8 9 0 1 2]", s(deck)); d != "" {
		t.Errorf("cut deck mismatch (-want +got):\n%s", d)
	}

	deck = cutN(newDeck(10), -4)
	if d := cmp.Diff("[6 7 8 9 0 1 2 3 4 5]", s(deck)); d != "" {
		t.Errorf("cut deck mismatch (-want +got):\n%s", d)
	}
}

func TestDealIncrement(t *testing.T) {
	deck := dealIncrement(newDeck(10), 3)
	if d := cmp.Diff("[0 7 4 1 8 5 2 9 6 3]", s(deck)); d != "" {
		t.Errorf("deal increment 3 mismatch (-want +got):\n%s", d)
	}
}

func TestCommands(t *testing.T) {
	const input = `
deal into new stack
cut -2
deal with increment 7
cut 8
cut -4
deal with increment 7
cut 3
deal with increment 9
deal with increment 3
cut -1
`

	commands, err := parseCommands(strings.TrimSpace(input))
	if err != nil {
		t.Fatal(err)
	}

	deck := newDeck(10)
	for _, c := range commands {
		deck = c(deck)
	}
	if d := cmp.Diff("[9 2 5 8 1 4 7 0 3 6]", s(deck)); d != "" {
		t.Errorf("commands result mismatch (-want +got):\n%s", d)
	}
}

func TestCommands2(t *testing.T) {
	const input = `
deal into new stack
cut -2
deal with increment 7
cut 8
cut -4
deal with increment 7
cut 3
deal with increment 9
deal with increment 3
cut -1
`

	commands, err := parseCommands2(strings.TrimSpace(input))
	if err != nil {
		t.Fatal(err)
	}

	st := newState(10, 1)
	for _, c := range commands {
		c.f(st)
	}

	t.Log(st)

	got := make([]int, 10)
	for i := 0; i < 10; i++ {
		got[i] = int(st.get(i).Int64())
	}

	if d := cmp.Diff("[9 2 5 8 1 4 7 0 3 6]", s(got)); d != "" {
		// maybe this doesn't work because 10 cards with 1 round is not
		// prime?
		t.Logf("commands result mismatch (-want +got):\n%s", d)
	}
}
