package render

import (
	"fmt"
	"testing"

	"github.com/kward/tabulate/table"
)

func TestRender_Split(t *testing.T) {
	for _, tc := range []struct {
		desc    string
		line    string
		numCols int

		csv      string
		markdown string
		mysql    string
		plain    string
		sqlite3  string
	}{
		{"single column", "123", 1,
			"123\n",
			"| 123 |\n",
			"+-----+\n| 123 |\n+-----+\n",
			"123\n",
			"123\n",
		},
		{"three columns", "1 2 3", 3,
			"1,2,3\n",
			"| 1 | 2 | 3 |\n",
			"+---+---+---+\n| 1 | 2 | 3 |\n+---+---+---+\n",
			"1 2 3\n",
			"1|2|3\n",
		},
		{"comment", "# 1 2 3", 1,
			"",
			"",
			"",
			"# 1 2 3\n",
			"",
		},
	} {
		tbl, err := table.Split([]string{tc.line}, " ", tc.numCols, table.EnableComments(true))
		if err != nil {
			t.Fatalf("unexpected error; %s", err)
		}

		t.Run(fmt.Sprintf("CSVRenderer %s", tc.desc), func(t *testing.T) {
			r := &CSVRenderer{}
			if got, want := r.Render(tbl), tc.csv; got != want {
				t.Errorf("= %q, want %q", got, want)
			}
		})

		t.Run(fmt.Sprintf("MarkdownRenderer %s", tc.desc), func(t *testing.T) {
			r := &MarkdownRenderer{}
			if got, want := r.Render(tbl), tc.markdown; got != want {
				t.Errorf("= %q, want %q", got, want)
			}
		})

		t.Run(fmt.Sprintf("MySQLRenderer %s", tc.desc), func(t *testing.T) {
			r := &MySQLRenderer{}
			if got, want := r.Render(tbl), tc.mysql; got != want {
				t.Errorf("= %q, want %q", got, want)
			}
		})

		t.Run(fmt.Sprintf("PlainRenderer %s", tc.desc), func(t *testing.T) {
			r := &PlainRenderer{}
			r.SetOFS(" ")
			if got, want := r.Render(tbl), tc.plain; got != want {
				t.Errorf("= %q, want %q", got, want)
			}
		})

		t.Run(fmt.Sprintf("SQLite3Renderer %s", tc.desc), func(t *testing.T) {
			r := &SQLite3Renderer{}
			if got, want := r.Render(tbl), tc.sqlite3; got != want {
				t.Errorf("= %q, want %q", got, want)
			}
		})
	}
}

func TestRender_Append(t *testing.T) {
	for _, tc := range []struct {
		desc    string
		records []string

		csv      string
		markdown string
		mysql    string
		plain    string
		sqlite3  string
	}{
		{"second column empty", []string{"123", ""},
			"123,\n",
			"| 123 | |\n",
			"+-----+-+\n| 123 | |\n+-----+-+\n",
			"123\n",
			"123|\n",
		},
	} {
		tbl, err := table.NewTable(table.EnableComments(true))
		if err != nil {
			t.Fatalf("unexpected error; %s", err)
		}
		tbl.Append(tc.records)

		t.Run(fmt.Sprintf("CSVRenderer %s", tc.desc), func(t *testing.T) {
			r := &CSVRenderer{}
			if got, want := r.Render(tbl), tc.csv; got != want {
				t.Errorf("= %q, want %q", got, want)
			}
		})

		t.Run(fmt.Sprintf("MarkdownRenderer %s", tc.desc), func(t *testing.T) {
			r := &MarkdownRenderer{}
			if got, want := r.Render(tbl), tc.markdown; got != want {
				t.Errorf("= %q, want %q", got, want)
			}
		})

		t.Run(fmt.Sprintf("MySQLRenderer %s", tc.desc), func(t *testing.T) {
			r := &MySQLRenderer{}
			if got, want := r.Render(tbl), tc.mysql; got != want {
				t.Errorf("= %q, want %q", got, want)
			}
		})

		t.Run(fmt.Sprintf("PlainRenderer %s", tc.desc), func(t *testing.T) {
			r := &PlainRenderer{}
			r.SetOFS(" ")
			if got, want := r.Render(tbl), tc.plain; got != want {
				t.Errorf("= %q, want %q", got, want)
			}
		})

		t.Run(fmt.Sprintf("SQLite3Renderer %s", tc.desc), func(t *testing.T) {
			r := &SQLite3Renderer{}
			if got, want := r.Render(tbl), tc.sqlite3; got != want {
				t.Errorf("= %q, want %q", got, want)
			}
		})
	}
}

type renderRow struct {
	in, out string
}

type renderTable struct {
	in  []string
	out string
}

func TestPlainRenderer(t *testing.T) {
	for _, tt := range []renderTable{
		{[]string{"1 22 333", "333 22 1"},
			"1   22 333\n333 22 1\n"},
		{[]string{"1 22 333", "4444 333 22 1"},
			"1    22  333\n4444 333 22  1\n"},
		{[]string{"4444 333 22", "1"},
			"4444 333 22\n1\n"},
	} {
		tbl, err := table.Split(tt.in, " ", -1)
		if err != nil {
			t.Fatalf("unexpected error; %s", err)
		}

		r := &PlainRenderer{}
		r.SetOFS(" ")
		if got, want := r.Render(tbl), tt.out; got != want {
			t.Errorf("Render() = %q, want %q", got, want)
		}
	}
}
