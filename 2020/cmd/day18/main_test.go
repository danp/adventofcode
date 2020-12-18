package main

import "testing"

func TestEvalPlus(t *testing.T) {
	cases := []struct {
		input string
		want  int
	}{
		{
			input: "1 + 2 * 3 + 4 * 5 + 6",
			want:  231,
		},
		{
			input: "1 + (2 * 3) + (4 * (5 + 6))",
			want:  51,
		},
		{
			input: "2 * 3 + (4 * 5)",
			want:  46,
		},
		{
			input: "5 + (8 * 3 + 9 + 3 * 4 * 3)",
			want:  1445,
		},
		{
			input: "5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))",
			want:  669060,
		},
		{
			input: "((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2",
			want:  23340,
		},
	}

	for _, tc := range cases {
		got := evalPlus(tc.input)
		if got != tc.want {
			t.Errorf("evalPlus(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}
