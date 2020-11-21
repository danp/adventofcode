package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExtractLayers(t *testing.T) {
	for _, tt := range []struct {
		input string
		w, h  int
		want  [][][]int
	}{
		{"123456789012", 3, 2, [][][]int{
			{{1, 2, 3}, {4, 5, 6}},
			{{7, 8, 9}, {0, 1, 2}},
		}},
	} {
		t.Run(tt.input, func(t *testing.T) {
			got, err := extractLayers(tt.input, tt.w, tt.h)
			if err != nil {
				t.Fatal(err)
			}

			if d := cmp.Diff(tt.want, got); d != "" {
				t.Errorf("layers mismatch (-want +got):\n%s", d)
			}
		})
	}
}

func TestFindFewestZeroes(t *testing.T) {
	const data = "0000000000000001"

	ls, err := extractLayers(data, 4, 2)
	if err != nil {
		t.Fatal(err)
	}

	l, err := findFewestZeroes(ls)
	if err != nil {
		t.Fatal(err)
	}

	if l != 1 {
		t.Errorf("got fewest zeroes layer %d, want %d", l, 1)
	}
}

func TestDecodeLayers(t *testing.T) {
	const data = "0222112222120000"

	layers, err := extractLayers(data, 2, 2)
	if err != nil {
		t.Fatal(err)
	}

	got, err := decodeLayers(layers)
	if err != nil {
		t.Fatal(err)
	}

	want := [][]int{
		{0, 1},
		{1, 0},
	}
	if d := cmp.Diff(want, got); d != "" {
		t.Errorf("decoded layer mismatch (-want +got):\n%s", d)
	}
}
