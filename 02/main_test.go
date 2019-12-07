package main

import "testing"

func TestPrograms(t *testing.T) {
	for _, tt := range []struct {
		input, output string
	}{
		{"1,9,10,3,2,3,11,0,99,30,40,50", "3500,9,10,70,2,3,11,0,99,30,40,50"},
		{"1,0,0,0,99", "2,0,0,0,99"},
		{"2,4,4,5,99,0", "2,4,4,5,99,9801"},
		{"1,1,1,4,99,5,6,0,99", "30,1,1,4,2,5,6,0,99"},
	} {
		t.Run(tt.input, func(t *testing.T) {
			p, err := parse(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			err = run(p)
			if err != nil {
				t.Fatal(err)
			}

			got := buildOut(p)
			if got != tt.output {
				t.Fatalf("run(%q) = %q, want %q", tt.input, got, tt.output)
			}
		})
	}
}
