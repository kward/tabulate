package tabulate

import (
	"strings"
	"testing"
)

const ifs = " "

type rowTests []struct {
	in, out string
}

type tableTests []struct {
	in  []string
	out string
}

func TestCSVRenderer(t *testing.T) {
	tests := rowTests{
		{"123", "123\n"},
		{"1 2 3", "1,2,3\n"},
		{"1  2  3", "1,2,3\n"},
		{"# 1 2 3", "# 1 2 3\n"},
		{"# 1  2  3", "# 1  2  3\n"},
	}
	testRowRendering(t, tests, &CSVRenderer{})
}

func TestMarkdownRenderer(t *testing.T) {
	tests := rowTests{
		{"123", "| 123 |\n"},
		{"1 2 3", "| 1 | 2 | 3 |\n"},
		{"1  2  3", "| 1 | 2 | 3 |\n"},
		{"# 1 2 3", ""},
	}
	testRowRendering(t, tests, &MarkdownRenderer{})
}

func TestMySQLRenderer(t *testing.T) {
	tests := rowTests{
		{"123", "+-----+\n| 123 |\n+-----+\n"},
		{"1 2 3", "+---+---+---+\n| 1 | 2 | 3 |\n+---+---+---+\n"},
		{"1  2  3", "+---+---+---+\n| 1 | 2 | 3 |\n+---+---+---+\n"},
		{"# 1 2 3", ""},
		{"# 1  2  3", ""},
	}
	testRowRendering(t, tests, &MySQLRenderer{})
}

func TestPlainRenderer(t *testing.T) {
	rt := rowTests{
		{"123", "123\n"},
		{"1 2 3", "1 2 3\n"},
		{"1  2  3", "1 2 3\n"},
		{"1  2  3", "1 2 3\n"},
		{"# 1 2 3", "# 1 2 3\n"},
		{"# 1  2  3", "# 1  2  3\n"},
	}
	testRowRendering(t, rt, &PlainRenderer{OFS: " "})

	tt := tableTests{
		{[]string{"1 22 333", "333 22 1"}, "1   22 333\n333 22 1\n"},
		{[]string{"1 22 333", "4444 333 22 1"}, "1    22  333\n4444 333 22  1\n"},
		{[]string{"4444 333 22", "1"}, "4444 333 22\n1\n"},
	}
	testTableRendering(t, tt, &PlainRenderer{OFS: " "})
}

func TestSQLite3Renderer(t *testing.T) {
	tests := rowTests{
		{"123", "123\n"},
		{"1 2 3", "1|2|3\n"},
		{"1  2  3", "1|2|3\n"},
		{"# 1 2 3", ""},
		{"# 1  2  3", ""},
	}
	testRowRendering(t, tests, &SQLite3Renderer{})
}

func testRowRendering(t *testing.T, tests rowTests, r Renderer) {
	tbl := NewTable(NewTableConfig())
	for _, tt := range tests {
		tbl.Split([]string{tt.in}, ifs, -1)
		if got, want := r.Render(&tbl), tt.out; got != want {
			got = strings.Replace(got, "\n", "\\n", -1)
			want = strings.Replace(want, "\n", "\\n", -1)
			t.Errorf("Render() = '%v', want '%v'", got, want)
		}
	}
}

func testTableRendering(t *testing.T, tests tableTests, r Renderer) {
	tbl := NewTable(NewTableConfig())
	for _, tt := range tests {
		tbl.Split(tt.in, ifs, -1)
		if got, want := r.Render(&tbl), tt.out; got != want {
			got = strings.Replace(got, "\n", "\\n", -1)
			want = strings.Replace(want, "\n", "\\n", -1)
			t.Errorf("Render() = '%v', want '%v'", got, want)
		}
	}
}
