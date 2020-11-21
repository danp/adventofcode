package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFFT(t *testing.T) {
	const input = "12345678"

	digits := convert(input)

	phase(digits)
	want := convert("48226158")
	if d := cmp.Diff(want, digits); d != "" {
		t.Errorf("phase 1 mismatch (-want +got):\n%s", d)
	}

	phase(digits)
	want = convert("34040438")
	if d := cmp.Diff(want, digits); d != "" {
		t.Errorf("phase 2 mismatch (-want +got):\n%s", d)
	}

	// larger
	digits = convert("80871224585914546619083218645595")
	for i := 0; i < 100; i++ {
		phase(digits)
	}
	want = convert("24176176")
	if d := cmp.Diff(want, digits[:8]); d != "" {
		t.Errorf("phase 1 mismatch (-want +got):\n%s", d)
	}
}
