package tabulate

import (
	"testing"

	"github.com/kward/golib/operators"
)

func TestSingleRowSplit(t *testing.T) {
	tests := []struct {
		records  []string
		columns  int
		colCount int
		colSizes []int
	}{
		{[]string{"1 2"}, -1, 2, []int{1, 1}},
		{[]string{"1 2 333"}, -1, 3, []int{1, 1, 3}},
		{[]string{"# comment"}, -1, 1, []int{0}},
		{[]string{"1 2 3"}, 2, 2, []int{1, 3}},
	}

	tbl := NewTable(NewTableConfig())
	for _, tt := range tests {
		tbl.Split(tt.records, " ", tt.columns)
		if tbl.colCount != tt.colCount {
			t.Errorf("Split(%v) colCount = %v, want %v", tt.records, tbl.colCount, tt.colCount)
		}
		if !operators.EqualSlicesOfInt(tbl.colSizes, tt.colSizes) {
			t.Errorf("Split(%v) colSizes = %v, want %v", tt.records, tbl.colSizes, tt.colSizes)
		}
	}
}

func TestMultiRowSplit(t *testing.T) {
	tests := []struct {
		records  []string
		columns  int
		colCount int
		colSizes []int
	}{
		{[]string{"1", "2 2"}, -1, 2, []int{1, 1}},
		{[]string{"1 22 333", "333 22 1"}, -1, 3, []int{3, 2, 3}},
	}

	tbl := NewTable(NewTableConfig())
	for _, tt := range tests {
		tbl.Split(tt.records, " ", tt.columns)
		if tbl.colCount != tt.colCount {
			t.Errorf("Split(%v) colCount = %v, want %v", tt.records, tbl.colCount, tt.colCount)
		}
		if !operators.EqualSlicesOfInt(tbl.colSizes, tt.colSizes) {
			t.Errorf("Split(%v) colSizes = %v, want %v", tt.records, tbl.colSizes, tt.colSizes)
		}
	}
}

func TestSplitRow(t *testing.T) {
	var (
		s string
		c []string
	)

	ifs := " "
	cols := 0
	tc := TableConfig{commentPrefix: "#", ignoreComments: true}

	s = "1 2 3"
	c = splitRow(s, ifs, cols, &tc)
	if !operators.EqualSlicesOfString(c, []string{"1", "2", "3"}) {
		t.Errorf("splitRow('%v', '%v', %v) => %v", s, ifs, cols, c)
	}

	s = "# comment line"
	if l := len(splitRow(s, ifs, cols, &tc)); l != 1 {
		t.Errorf("splitRow('%v', '%v', %v): len %v", s, ifs, cols, l)
	}
}

func TestIsComment(t *testing.T) {
	var tests = []struct {
		in  []string
		out bool
	}{
		{in: []string{"foo"}, out: false},
		{in: []string{"foo", "bar"}, out: false},
		{in: []string{"# foo"}, out: true},
		{in: []string{"# foo", "bar"}, out: false},
	}

	tbl := NewTable(NewTableConfig())
	for _, tt := range tests {
		if got, want := tbl.IsComment(tt.in), tt.out; got != want {
			t.Errorf("IsComment() = %v, want %v", got, want)
		}
	}
}

func TestSplitNMerged(t *testing.T) {
	t.Parallel()

	var got, want []string

	want = []string{"1", "2", "3"}
	got = SplitNMerged("1 2 3", " ", -1)
	if !operators.EqualSlicesOfString(want, got) {
		t.Errorf("SplitNMerged(): want %v, got %v", want, got)
	}

	want = []string{"1", "2", "3"}
	got = SplitNMerged("1 2   3", " ", -1)
	if !operators.EqualSlicesOfString(want, got) {
		t.Errorf("SplitNMerged(): want %v, got %v", want, got)
	}

	want = []string{}
	got = SplitNMerged("", " ", -1)
	if !operators.EqualSlicesOfString(want, got) {
		t.Errorf("SplitNMerged(): want %v, got %v", want, got)
	}
}

func TestStretch(t *testing.T) {
	t.Parallel()

	var got, have, want string

	// Long.
	have = "str1"
	got = Stretch(have, ' ', 6)
	want = "str1  "
	if want != got {
		t.Errorf("Stretch(%q): got %v, want %v", have, want, got)
	}
	// Just right.
	have = "str2"
	got = Stretch(have, ' ', 4)
	want = "str2"
	if want != got {
		t.Errorf("Stretch(%q): got %v, want %v", have, want, got)
	}
	// Short.
	have = "str3"
	got = Stretch(have, ' ', 1)
	want = "str3"
	if want != got {
		t.Errorf("Stretch(%q): got %v, want %v", have, want, got)
	}
}
