package main

import (
	"reflect"
	"testing"
)

func TestGenerateTomography(t *testing.T) {
	testCases := []struct {
		name         string
		state        string
		wantRowHints [][]int
		wantColHints [][]int
	}{
		{
			name:         "empty state",
			state:        "",
			wantRowHints: [][]int{{0}},
			wantColHints: [][]int{},
		},
		{
			name:         "single cell filled",
			state:        "1",
			wantRowHints: [][]int{{1}},
			wantColHints: [][]int{{1}},
		},
		{
			name:         "single cell empty",
			state:        " ",
			wantRowHints: [][]int{{0}},
			wantColHints: [][]int{{0}},
		},
		{
			name: "simple 2x2",
			state: "1 \n" +
				" 1",
			wantRowHints: [][]int{{1}, {1}},
			wantColHints: [][]int{{1}, {1}},
		},
		{
			name: "complex rows",
			state: "11 1\n" +
				"1 11",
			wantRowHints: [][]int{{2, 1}, {1, 2}},
			wantColHints: [][]int{{2}, {1}, {1}, {2}},
		},
		{
			name: "all filled",
			state: "11\n" +
				"11",
			wantRowHints: [][]int{{2}, {2}},
			wantColHints: [][]int{{2}, {2}},
		},
		{
			name: "all empty",
			state: "  \n" +
				"  ",
			wantRowHints: [][]int{{0}, {0}},
			wantColHints: [][]int{{0}, {0}},
		},
		{
			name: "empty last column",
			state: "1 \n" +
				"1 \n" +
				"1 ",
			wantRowHints: [][]int{{1}, {1}, {1}},
			wantColHints: [][]int{{3}, {0}},
		},
		{
			name: "empty first column",
			state: " 1\n" +
				" 1\n" +
				" 1",
			wantRowHints: [][]int{{1}, {1}, {1}},
			wantColHints: [][]int{{0}, {3}},
		},
		{
			name: "empty last row",
			state: "11\n" +
				"  ",
			wantRowHints: [][]int{{2}, {0}},
			wantColHints: [][]int{{1}, {1}},
		},
		{
			name: "empty first row",
			state: "  \n" +
				"11",
			wantRowHints: [][]int{{0}, {2}},
			wantColHints: [][]int{{1}, {1}},
		},
		{
			name:         "jagged input", // width is determined by first line
			state:        "111\n11\n1",
			wantRowHints: [][]int{{3}, {2}, {1}},
			wantColHints: [][]int{{3}, {2}, {1}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotRowHints, gotColHints := generateTomography(tc.state)

			if !reflect.DeepEqual(gotRowHints, tc.wantRowHints) {
				t.Errorf("generateTomography() row hints = %v, want %v", gotRowHints, tc.wantRowHints)
			}
			if !reflect.DeepEqual(gotColHints, tc.wantColHints) {
				t.Errorf("generateTomography() col hints = %v, want %v", gotColHints, tc.wantColHints)
			}
		})
	}
}
