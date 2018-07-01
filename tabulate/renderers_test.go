package tabulate

import (
	"testing"
)

type renderRow struct {
	in, out string
}

type renderTable struct {
	in  []string
	out string
}

func TestCSVRenderer(t *testing.T) {
	tbl, err := NewTable()
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}
	r := &CSVRenderer{}

	for _, tt := range []renderRow{
		{"123", "123\n"},
		{"1 2 3", "1,2,3\n"},
		{"1  2  3", "1,2,3\n"},
		{"# 1 2 3", "# 1 2 3\n"},
		{"# 1  2  3", "# 1  2  3\n"},
	} {
		tbl.Split([]string{tt.in}, " ", -1)
		if got, want := r.Render(tbl), tt.out; got != want {
			t.Errorf("Render() = %q, want %q", got, want)
		}
	}
}

func TestMarkdownRenderer(t *testing.T) {
	tbl, err := NewTable()
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}
	r := &MarkdownRenderer{}

	for _, tt := range []renderRow{
		{"123", "| 123 |\n"},
		{"1 2 3", "| 1 | 2 | 3 |\n"},
		{"1  2  3", "| 1 | 2 | 3 |\n"},
		{"# 1 2 3", ""},
	} {
		tbl.Split([]string{tt.in}, " ", -1)
		if got, want := r.Render(tbl), tt.out; got != want {
			t.Errorf("Render() = %q, want %q", got, want)
		}
	}
}

func TestMySQLRenderer(t *testing.T) {
	tbl, err := NewTable()
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}
	r := &MySQLRenderer{}

	for _, tt := range []renderRow{
		{"123", "+-----+\n| 123 |\n+-----+\n"},
		{"1 2 3", "+---+---+---+\n| 1 | 2 | 3 |\n+---+---+---+\n"},
		{"1  2  3", "+---+---+---+\n| 1 | 2 | 3 |\n+---+---+---+\n"},
		{"# 1 2 3", ""},
		{"# 1  2  3", ""},
	} {
		tbl.Split([]string{tt.in}, " ", -1)
		if got, want := r.Render(tbl), tt.out; got != want {
			t.Errorf("Render() = %q, want %q", got, want)
		}
	}
}

func TestPlainRenderer(t *testing.T) {
	tbl, err := NewTable()
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}
	r := &PlainRenderer{}
	r.SetOFS(" ")

	for _, tt := range []renderRow{
		{"123", "123\n"},
		{"1 2 3", "1 2 3\n"},
		{"1  2  3", "1 2 3\n"},
		{"1  2  3", "1 2 3\n"},
		{"# 1 2 3", "# 1 2 3\n"},
		{"# 1  2  3", "# 1  2  3\n"},
	} {
		tbl.Split([]string{tt.in}, " ", -1)
		if got, want := r.Render(tbl), tt.out; got != want {
			t.Errorf("Render() = %q, want %q", got, want)
		}
	}

	for _, tt := range []renderTable{
		{[]string{"1 22 333", "333 22 1"}, "1   22 333\n333 22 1\n"},
		{[]string{"1 22 333", "4444 333 22 1"}, "1    22  333\n4444 333 22  1\n"},
		{[]string{"4444 333 22", "1"}, "4444 333 22\n1\n"},
	} {
		tbl.Split(tt.in, " ", -1)
		if got, want := r.Render(tbl), tt.out; got != want {
			t.Errorf("Render() = %q, want %q", got, want)
		}
	}
}

func TestSQLite3Renderer(t *testing.T) {
	tbl, err := NewTable()
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}
	r := &SQLite3Renderer{}

	for _, tt := range []renderRow{
		{"123", "123\n"},
		{"1 2 3", "1|2|3\n"},
		{"1  2  3", "1|2|3\n"},
		{"# 1 2 3", ""},
		{"# 1  2  3", ""},
	} {
		tbl.Split([]string{tt.in}, " ", -1)
		if got, want := r.Render(tbl), tt.out; got != want {
			t.Errorf("Render() = %q, want %q", got, want)
		}
	}
}
