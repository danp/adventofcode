package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var totalFuel int

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		mass, err := strconv.Atoi(strings.TrimSpace(sc.Text()))
		if err != nil {
			panic(err)
		}

		totalFuel += fuelCounterUpperWithFuel(mass)
	}

	fmt.Println(totalFuel)
}

func fuelCounterUpper(mass int) int {
	fuel := mass / 3
	fuel -= 2
	if fuel < 0 {
		fuel = 0
	}
	return fuel
}

func fuelCounterUpperWithFuel(mass int) int {
	var out int

	for mass > 0 {
		fuelMass := fuelCounterUpper(mass)
		out += fuelMass
		mass = fuelMass
	}

	return out
}
