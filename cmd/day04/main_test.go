package main

import (
	"strconv"
	"testing"
)

func TestCheckPassword(t *testing.T) {
	for _, tt := range []struct {
		in   int
		want bool
	}{
		{122345, true},
		{111123, true},
		{135679, false}, // no double
		{111111, true},
		{111112, true},
		{111121, false},
		{911111, false},
		{223450, false},
		{123789, false},
		{888889, true},
	} {
		t.Run(strconv.Itoa(tt.in), func(t *testing.T) {
			if got := checkPassword(tt.in); got != tt.want {
				t.Fatalf("got result %t, want %t", got, tt.want)
			}
		})
	}
}

func TestCheckPasswordHarder(t *testing.T) {
	for _, tt := range []struct {
		in   int
		want bool
	}{
		{122345, true},
		{111123, false},
		{135679, false}, // no double
		{111111, false},
		{111112, false},
		{111121, false},
		{911111, false},
		{223450, false},
		{123789, false},
		{888889, false},
		{112233, true},
		{123444, false},
		{111122, true},
		{223344, true},
		{555566, true},
		{123345, true},
		{112345, true},
	} {
		t.Run(strconv.Itoa(tt.in), func(t *testing.T) {
			if got := checkPasswordHarder(tt.in); got != tt.want {
				t.Fatalf("got result %t, want %t", got, tt.want)
			}
		})
	}
}
