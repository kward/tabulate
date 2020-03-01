package table

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kward/golib/operators"
)

func TestAppend(t *testing.T) {
	for _, tc := range []struct {
		desc  string
		elems [][]string
		row   *Row
	}{
		{"a b c",
			[][]string{{"a", "b", "c"}},
			&Row{
				[]*Column{&Column{"a"}, &Column{"b"}, &Column{"c"}},
				[]int{1, 1, 1},
				false},
		},
		{"empty",
			[][]string{},
			&Row{
				[]*Column{},
				[]int{},
				false}},
	} {
		t.Run(fmt.Sprintf("Append() %s", tc.desc), func(t *testing.T) {
			tbl, err := NewTable()
			if err != nil {
				t.Fatalf("unexpected error; %s", err)
			}

			tbl.Append(tc.elems...)
			if got, want := len(tbl.rows), len(tc.elems); got != want {
				t.Errorf("len(rows) = %v, want %v", got, want)
			}
			if len(tc.elems) == 0 {
				return
			}
			if got, want := tbl.Rows()[0].Values(), tc.row.Values(); !reflect.DeepEqual(got, want) {
				t.Errorf("row = %q, want %q", got, want)
			}
		})
	}
}

func TestSplit_SingleRow(t *testing.T) {
	for _, tc := range []struct {
		desc   string
		lines  []string
		splits int

		numRows   int
		numCols   int
		values    []string
		sizes     []int // Slice of column sizes.
		isComment bool
	}{
		// Normal cases.
		{"one column", []string{"1"}, 1,
			1, 1, []string{"1"}, []int{1}, false},
		{"two columns", []string{"1 2"}, 2,
			1, 2, []string{"1", "2"}, []int{1, 1}, false},
		{"two>one column", []string{"1 2"}, 1,
			1, 1, []string{"1 2"}, []int{3}, false},
		{"three seen columns", []string{"1 2 333"}, -1,
			1, 3, []string{"1", "2", "333"}, []int{1, 1, 3}, false},
		{"three seen wide columns", []string{"1  2  3"}, -1,
			1, 3, []string{"1", "2", "3"}, []int{1, 1, 1}, false},

		// Comments.
		{"comment", []string{"# comment"}, -1,
			1, 1, []string{"# comment"}, []int{9}, true},

		// Special cases.
		{desc: "zero splits", lines: []string{"abc 123"}},
	} {
		t.Run(fmt.Sprintf("Split() single-row %s", tc.desc), func(t *testing.T) {
			tbl, err := Split(tc.lines, " ", tc.splits, EnableComments(true))
			if err != nil {
				t.Fatalf("unexpected error; %s", err)
			}

			if got, want := tbl.NumRows(), tc.numRows; got != want {
				t.Errorf("tbl.NumRows() = %d, want %d", got, want)
			}
			if tc.splits == 0 {
				return
			}

			row := tbl.Rows()[0]
			if got, want := row.NumColumns(), tc.numCols; got != want {
				t.Errorf("row.NumColumns() = %d, want %d", got, want)
			}
			if got, want := row.Values(), tc.values; !operators.EqualSlicesOfString(got, want) {
				t.Errorf("row.Values() = %q, want %q", got, want)
			}
			if got, want := row.Sizes(), tc.sizes; !operators.EqualSlicesOfInt(got, want) {
				t.Errorf("row.Sizes() = %d, want %d", got, want)
			}
		})
	}
}

func TestParse_MultiRow(t *testing.T) {
	for _, tc := range []struct {
		desc   string
		lines  []string
		splits int

		numRows  int
		colSizes []int

		numCols   []int
		isComment []bool
	}{
		// Normal cases.
		{"2x2", []string{"1", "2 2"}, -1,
			2, []int{1, 1}, []int{1, 2}, []bool{false, false}},
		{"2x3", []string{"1 22", "22 333", "333 4444"}, -1,
			3, []int{3, 4}, []int{2, 2, 2}, []bool{false, false, false}},
		{"3x2", []string{"1 22 333", "333 22 1"}, -1,
			2, []int{3, 2, 3}, []int{3, 3}, []bool{false, false, false}},

		// Comment cases.
		{"2x3 with comment", []string{"1 22", "22 333", "# 333 4444"}, -1,
			3, []int{2, 3}, []int{2, 2, 1}, []bool{false, false, true}},

		// Special cases.
		{desc: "zero splits", lines: []string{"abc def", "123 456"}},
	} {
		t.Run(fmt.Sprintf("Split() multi-row %s", tc.desc), func(t *testing.T) {
			tbl, err := Split(tc.lines, " ", tc.splits, EnableComments(true))
			if err != nil {
				t.Fatalf("unexpected error; %s", err)
			}

			if got, want := tbl.NumRows(), tc.numRows; got != want {
				t.Errorf("tbl.NumRows() = %d, want %d", got, want)
			}
			if tc.splits == 0 {
				return
			}
			for i, s := range tbl.ColSizes() {
				if got, want := s, tc.colSizes[i]; got != want {
					t.Errorf("column #%d: tbl.ColSizes() = %d, want %d", i, got, want)
				}
			}

			for i, _ := range tc.lines {
				row := tbl.Rows()[i]
				if got, want := row.NumColumns(), tc.numCols[i]; got != want {
					t.Errorf("row #%d: row.NumColumns() = %d, want %d", i, got, want)
				}
				if got, want := row.IsComment(), tc.isComment[i]; got != want {
					t.Errorf("row #%d: row.IsComment() = %t", i, got)
				}
			}
		})
	}
}
