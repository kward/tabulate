package tabulate

import (
	"fmt"
	"testing"

	"github.com/kward/golib/operators"
)

func TestSplit_SingleRow(t *testing.T) {
	tbl, err := NewTable()
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}

	for _, tc := range []struct {
		desc    string
		lines   []string
		splits  int
		columns int   // Count of columns.
		sizes   []int // Slice of column sizes.
	}{
		{"one column", []string{"1"}, 1, 1, []int{1}},
		{"two columns", []string{"1 2"}, 2, 2, []int{1, 1}},
		{"two>one column", []string{"1 2"}, -1, 2, []int{1}},
		{"three columns", []string{"1 2 333"}, -1, 3, []int{1, 1, 3}},
		{"comment", []string{"# comment"}, -1, 1, []int{0}},
	} {
		t.Run(fmt.Sprintf("Split() single-row %s", tc.desc), func(t *testing.T) {
			tbl.Split(tc.lines, " ", tc.splits)
			row := tbl.rows[0]
			if got, want := row.Columns(), tc.columns; got != want {
				t.Errorf("%s: row.Columns() = %d, want %d", tc.desc, got, want)
			}
		})
	}
}

func TestSplit_MultiRow(t *testing.T) {
	tbl, err := NewTable()
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}

	for _, tc := range []struct {
		desc    string
		lines   []string
		splits  int
		columns []int
		sizes   []int
	}{
		{"2x2", []string{"1", "2 2"}, -1, []int{1, 2}, []int{1, 1}},
		{"3x2", []string{"1 22 333", "333 22 1"}, -1, []int{3, 3}, []int{3, 2, 3}},
	} {
		t.Run(fmt.Sprintf("Split() multi-row %s", tc.desc), func(t *testing.T) {
			tbl.Split(tc.lines, " ", tc.splits)
			for i := 0; i < len(tc.lines); i++ {
				row := tbl.rows[i]
				if got, want := row.Columns(), tc.columns[i]; got != want {
					t.Errorf("columns = %d, want %d", got, want)
				}
			}
		})
	}
}

func TestSplitLine(t *testing.T) {
	tbl, err := NewTable()
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}

	for _, tc := range []struct {
		desc  string
		line  string
		ifs   string
		cols  int
		split []string
	}{
		{"auto cols narrow", "1 2 3", " ", 0, []string{"1", "2", "3"}},
		{"auto cols wide", "1   2   3", " ", 0, []string{"1", "2", "3"}},
		{"one cols narrow", "1 2 3", " ", 1, []string{"1 2 3"}},
		{"two cols narrow", "1 2 3", " ", 2, []string{"1", "2 3"}},
		{"three cols narrow", "1 2 3", " ", 3, []string{"1", "2", "3"}},
		{"comment", "# comment line", " ", 0, []string{"# comment line"}},
	} {
		t.Run(fmt.Sprintf("splitLine() %s", tc.desc), func(t *testing.T) {
			got, want := tbl.splitLine(tc.line, tc.ifs, tc.cols), tc.split
			if !operators.EqualSlicesOfString(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestIsComment(t *testing.T) {
	tbl, err := NewTable()
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}

	for _, tc := range []struct {
		records   []string
		isComment bool
	}{
		{[]string{"foo"}, false},
		{[]string{"foo", "bar"}, false},
		{[]string{"# foo"}, true},
		{[]string{"# foo", "bar"}, false},
	} {
		row := NewRow(tc.records)
		if got, want := tbl.IsComment(row), tc.isComment; got != want {
			t.Errorf("IsComment(%v) = %v, want %v", tc.records, got, want)
		}
	}
}
