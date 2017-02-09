package tabulate

import (
	"strings"

	mymath "github.com/kward/tabulate/go/math"
	mystrings "github.com/kward/tabulate/go/strings"
)

type TableConfig struct {
	commentPrefix  string
	ignoreComments bool
}

func NewTableConfig() TableConfig {
	return TableConfig{"#", true}
}

const MAX_COLS = 100

type Table struct {
	records  [][]string
	colCount int   // Total number of columns.
	colSizes []int // The maximum size of each column.
	config   TableConfig
}

func NewTable(config TableConfig) Table {
	return Table{config: config}
}

func splitRow(row string, ifs string, columns int, config *TableConfig) []string {
	if config.ignoreComments && strings.HasPrefix(row, config.commentPrefix) {
		return []string{row}
	}

	if columns == 0 {
		columns = -1
	}
	return mystrings.SplitNMerged(row, ifs, columns)
}

func (t *Table) Split(records []string, ifs string, columns int) {
	var (
		cols    []string
		rows    [][]string
		numCols int
	)

	colSizes := make([]int, 0, MAX_COLS)

	for _, row := range records {
		cols = splitRow(row, ifs, columns, &t.config)
		if len(cols) == 1 && strings.HasPrefix(cols[0], t.config.commentPrefix) {
			if numCols == 0 {
				colSizes = append(colSizes, 0)
				numCols++
			}
		} else {
			for i, col := range cols {
				if i == numCols {
					colSizes = append(colSizes, 0)
					numCols++
				}
				colSizes[i] = mymath.Max(len(col), colSizes[i])
			}
		}
		rows = append(rows, cols)
	}

	t.records = rows
	t.colCount = len(colSizes)
	t.colSizes = colSizes
}

func (t *Table) IsComment(row []string) bool {
	return len(row) == 1 && strings.HasPrefix(row[0], t.config.commentPrefix)
}

// SplitNMerged slices s into substrings separated by sep and returns a slice
// of the substrings between those separators. If sep is empty, SplitN splits
// after each UTF-8 sequence. Repeated sep will be merged. The count determines
// the number of substrings to return:
//
//  n > 0: at most n substrings; the last substring will be the unsplit remainder.
//  n == 0: the result is nil (zero substrings)
//  n < 0: all substrings
func SplitNMerged(s string, sep string, n int) []string {
	split := strings.SplitN(s, sep, n)
	merged := make([]string, 0, len(split))
	for _, v := range split {
		if len(v) > 0 {
			merged = append(merged, v)
		}
	}
	return merged
}

// Stretches a string to a given length by appending a character to the end.
func Stretch(s string, r rune, l int) string {
	// Special cases.
	if len(s) >= l {
		return s
	}
	return s + strings.Repeat(string(r), l-len(s))
}
