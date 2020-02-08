package tabulate

import (
	"fmt"
	"testing"
)

type renderRow struct {
	in, out string
}

type renderTable struct {
	in  []string
	out string
}

func TestRender(t *testing.T) {
	tbl, err := NewTable()
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}

	for _, tc := range []struct {
		in       string
		cols     int
		csv      string
		markdown string
		mysql    string
		plain    string
		sqlite3  string
	}{
		{"123", 1,
			"123\n",
			"| 123 |\n",
			"+-----+\n| 123 |\n+-----+\n",
			"123\n",
			"123\n",
		},
		{"1 2 3", 3,
			"1,2,3\n",
			"| 1 | 2 | 3 |\n",
			"+---+---+---+\n| 1 | 2 | 3 |\n+---+---+---+\n",
			"1 2 3\n",
			"1|2|3\n",
		},
		{"1  2  3", 3,
			"1,2,3\n",
			"| 1 | 2 | 3 |\n",
			"+---+---+---+\n| 1 | 2 | 3 |\n+---+---+---+\n",
			"1 2 3\n",
			"1|2|3\n",
		},
		{"# 1 2 3", 1,
			"",
			"",
			"",
			"# 1 2 3\n",
			"",
		},
		{"# 1  2  3", 1,
			"",
			"",
			"",
			"# 1  2  3\n",
			"",
		},
	} {
		t.Run(fmt.Sprintf("CSVRenderer"), func(t *testing.T) {
			r := &CSVRenderer{}
			tbl.Split([]string{tc.in}, " ", tc.cols)
			if got, want := r.Render(tbl), tc.csv; got != want {
				t.Errorf("Render(%q) = %q, want %q", tc.in, got, want)
			}
		})

		t.Run(fmt.Sprintf("MarkdownRenderer"), func(t *testing.T) {
			r := &MarkdownRenderer{}
			tbl.Split([]string{tc.in}, " ", tc.cols)
			if got, want := r.Render(tbl), tc.markdown; got != want {
				t.Errorf("Render(%q) = %q, want %q", tc.in, got, want)
			}
		})

		t.Run(fmt.Sprintf("MarkdownRenderer"), func(t *testing.T) {
			r := &MySQLRenderer{}
			tbl.Split([]string{tc.in}, " ", tc.cols)
			if got, want := r.Render(tbl), tc.mysql; got != want {
				t.Errorf("Render(%q) = %q, want %q", tc.in, got, want)
			}
		})

		t.Run(fmt.Sprintf("PlainRenderer"), func(t *testing.T) {
			r := &PlainRenderer{}
			r.SetOFS(" ")
			tbl.Split([]string{tc.in}, " ", tc.cols)
			if got, want := r.Render(tbl), tc.plain; got != want {
				t.Errorf("Render(%q) = %q, want %q", tc.in, got, want)
			}
		})

		t.Run(fmt.Sprintf("SQLite3Renderer"), func(t *testing.T) {
			r := &SQLite3Renderer{}
			tbl.Split([]string{tc.in}, " ", tc.cols)
			if got, want := r.Render(tbl), tc.sqlite3; got != want {
				t.Errorf("Render(%q) = %q, want %q", tc.in, got, want)
			}
		})
	}
}

func TestPlainRenderer(t *testing.T) {
	tbl, err := NewTable()
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}
	r := &PlainRenderer{}
	r.SetOFS(" ")

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
