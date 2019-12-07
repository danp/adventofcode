package main

import "testing"

func TestOrbitCount(t *testing.T) {
	const data = `
COM)B
B)C
C)D
D)E
E)F
B)G
G)H
D)I
E)J
J)K
K)L
`
	d, err := parse(data)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := orbitCount(d), 42; got != want {
		t.Fatalf("got %d direct and indirect orbits, want %d", got, want)
	}
}

func TestTransferCount(t *testing.T) {
	const data = `
COM)B
B)C
C)D
D)E
E)F
B)G
G)H
D)I
E)J
J)K
K)L
K)YOU
I)SAN
`

	d, err := parse(data)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := transferCount(d, "YOU", "SAN"), 4; got != want {
		t.Fatalf("got %d transfers, want %d", got, want)
	}
}
