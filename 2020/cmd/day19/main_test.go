package main

import (
	"reflect"
	"testing"
)

func TestRuleMatch(t *testing.T) {
	cases := []struct {
		start string
		input string
		want  []string
	}{
		{
			start: "0",
			input: "ab",
			want:  []string{"ab"},
		},
		{
			start: "0",
			input: "abb",
			want:  []string{"ab", "abb"},
		},
		{
			start: "0",
			input: "abbb",
			want:  []string{"ab", "abb", "abbb"},
		},
		{
			start: "3",
			input: "a",
			want:  []string{"a"},
		},
		{
			start: "3",
			input: "b",
			want:  nil,
		},
		{
			start: "3",
			input: "a",
			want:  []string{"a"},
		},
		{
			start: "1",
			input: "ab",
			want:  []string{"ab"},
		},
		{
			start: "1",
			input: "ba",
			want:  []string{"ba"},
		},
		{
			start: "5",
			input: "aba",
			want:  []string{"aba"},
		},
		{
			start: "5",
			input: "abbbba",
			want:  []string{"abbbba"},
		},
		{
			start: "5",
			input: "abbbbbbbbbbbba",
			want:  []string{"abbbbbbbbbbbba"},
		},
	}

	rs := rules{
		g: map[string]rule{
			"0": {name: "0", seqs: [][]string{{"3", "2"}}},
			"1": {name: "1", seqs: [][]string{{"3", "4"}, {"4", "3"}}},
			"2": {name: "2", seqs: [][]string{{"4"}, {"4", "2"}}},
			"3": {name: "3", lit: "a"},
			"4": {name: "4", lit: "b"},
			"5": {name: "5", seqs: [][]string{{"3", "4", "3"}, {"3", "4", "2", "3"}}},
		},
	}

	for _, tc := range cases {
		got := rs.g[tc.start].match(rs.g, tc.input)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("rule[%q].match(%q) = %+v, want %+v", tc.start, tc.input, got, tc.want)
		}
	}
}
