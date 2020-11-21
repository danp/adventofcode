package intcode

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPrograms(t *testing.T) {
	for _, tt := range []struct {
		input                         string
		memSize                       int
		output                        string
		providedInput, expectedOutput []int
	}{
		{"1,9,10,3,2,3,11,0,99,30,40,50", 0, "3500,9,10,70,2,3,11,0,99,30,40,50", nil, nil},
		{"1,0,0,0,99", 0, "2,0,0,0,99", nil, nil},
		{"2,4,4,5,99,0", 0, "2,4,4,5,99,9801", nil, nil},
		{"1,1,1,4,99,5,6,0,99", 0, "30,1,1,4,2,5,6,0,99", nil, nil},
		{"1002,4,3,4,33", 0, "1002,4,3,4,99", nil, nil},
		{"3,3,99,123", 0, "3,3,99,456", []int{456}, nil},
		{"4,3,99,123", 0, "4,3,99,123", nil, []int{123}},
		{"1101,100,-1,4,0", 0, "1101,100,-1,4,99", nil, nil},
		{"5,7,4,99,4,7,99,1", 0, "5,7,4,99,4,7,99,1", nil, []int{1}},
		{"6,7,4,99,4,7,99,0", 0, "6,7,4,99,4,7,99,0", nil, []int{0}},
		{"1107,1,2,5,99,0", 0, "1107,1,2,5,99,1", nil, nil},
		{"1107,2,1,5,99,1", 0, "1107,2,1,5,99,0", nil, nil},
		{"1108,1,1,5,99,0", 0, "1108,1,1,5,99,1", nil, nil},
		{"1108,2,1,5,99,1", 0, "1108,2,1,5,99,0", nil, nil},
		{"3,9,8,9,10,9,4,9,99,-1,8", 0, "3,9,8,9,10,9,4,9,99,1,8", []int{8}, []int{1}},
		{"3,9,8,9,10,9,4,9,99,-1,8", 0, "3,9,8,9,10,9,4,9,99,0,8", []int{7}, []int{0}},
		{"3,9,7,9,10,9,4,9,99,-1,8", 0, "3,9,7,9,10,9,4,9,99,0,8", []int{8}, []int{0}},
		{"3,9,7,9,10,9,4,9,99,-1,8", 0, "3,9,7,9,10,9,4,9,99,1,8", []int{7}, []int{1}},
		{"3,3,1108,-1,8,3,4,3,99", 0, "3,3,1108,1,8,3,4,3,99", []int{8}, []int{1}},
		{"3,3,1108,-1,8,3,4,3,99", 0, "3,3,1108,0,8,3,4,3,99", []int{7}, []int{0}},
		{"109,3,204,2,99,42", 0, "109,3,204,2,99,42", nil, []int{42}},
		{"1102,34915192,34915192,7,4,7,99,0", 0, "1102,34915192,34915192,7,4,7,99,1219070632396864", nil, []int{1219070632396864}},
		{"104,1125899906842624,99", 0, "104,1125899906842624,99", nil, []int{1125899906842624}},
		{"109,1,204,-1,1001,23,1,23,1008,23,16,24,1006,24,0,99", 25, "109,1,204,-1,1001,23,1,23,1008,23,16,24,1006,24,0,99,0,0,0,0,0,0,0,16,1", nil, []int{109, 1, 204, -1, 1001, 23, 1, 23, 1008, 23, 16, 24, 1006, 24, 0, 99}},
		{"109,5,203,2,99,0,0,0", 0, "109,5,203,2,99,0,0,42", []int{42}, nil},
	} {
		t.Run(tt.input, func(t *testing.T) {
			p, err := Parse(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			input := func() (int, error) {
				if len(tt.providedInput) == 0 {
					t.Error("input called with no more provided input")
				}
				in := tt.providedInput[0]
				tt.providedInput = tt.providedInput[1:]
				return in, nil
			}

			var gotOutput []int
			output := func(x int) error {
				gotOutput = append(gotOutput, x)
				return nil
			}

			memSize := tt.memSize
			if memSize == 0 {
				memSize = len(p)
			}
			mem := make([]int, memSize)
			err = Run(p, mem, input, output)
			if err != nil {
				t.Fatal(err)
			}

			got := buildOut(mem)
			if got != tt.output {
				t.Errorf("run(%q) = %q, want %q", tt.input, got, tt.output)
			}

			if d := cmp.Diff(tt.expectedOutput, gotOutput); d != "" {
				t.Errorf("output mismatch (-want +got):\n%s", d)
			}
		})
	}
}

func TestAProgram(t *testing.T) {
	const prog = "3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99"

	for _, tt := range []struct {
		input, want int
	}{
		{7, 999},
		{8, 1000},
		{9, 1001},
	} {
		t.Run(strconv.Itoa(tt.input), func(t *testing.T) {
			p, err := Parse(prog)
			if err != nil {
				t.Fatal(err)
			}

			input := func() (int, error) {
				return tt.input, nil
			}

			var got int
			output := func(x int) error {
				got = x
				return nil
			}

			mem := make([]int, 1024)
			if err := Run(p, mem, input, output); err != nil {
				t.Fatal(err)
			}

			if got != tt.want {
				t.Fatalf("got output %d, want %d", got, tt.want)
			}
		})
	}
}
