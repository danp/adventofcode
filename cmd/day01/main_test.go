package main

import "testing"

func TestFuelCounterUpper(t *testing.T) {
	for _, tt := range []struct {
		mass, fuel int
	}{
		{12, 2},
		{14, 2},
		{1969, 654},
		{100756, 33583},
	} {
		got := fuelCounterUpper(tt.mass)
		if got != tt.fuel {
			t.Errorf("fuelCounterUpper(%d) = %d, want %d", tt.mass, got, tt.fuel)
		}
	}
}

func TestFuelCounterUpperWithFuel(t *testing.T) {
	for _, tt := range []struct {
		mass, fuel int
	}{
		{14, 2},
		{1969, 966},
		{100756, 50346},
	} {
		got := fuelCounterUpperWithFuel(tt.mass)
		if got != tt.fuel {
			t.Errorf("fuelCounterUpperWithFuel(%d) = %d, want %d", tt.mass, got, tt.fuel)
		}
	}

}
