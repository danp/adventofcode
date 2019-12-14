package main

import (
	"testing"
)

func TestCountOre(t *testing.T) {
	reactions := map[string]reaction{
		"FUEL": reaction{1, "FUEL", map[string]int{"A": 7, "E": 1}},
		"E":    reaction{1, "E", map[string]int{"A": 7, "D": 1}},
		"D":    reaction{1, "D", map[string]int{"A": 7, "C": 1}},
		"C":    reaction{1, "C", map[string]int{"A": 7, "B": 1}},
		"B":    reaction{1, "B", map[string]int{"ORE": 1}},
		"A":    reaction{10, "A", map[string]int{"ORE": 10}},
		"ORE":  reaction{1, "ORE", nil},
	}

	surplus := make(map[string]int)
	got := countOre(reactions, "FUEL", 1, surplus)

	if got != 31 {
		t.Errorf("got ore %d, want %d", got, 31)
	}
}
