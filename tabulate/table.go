package tabulate

import (
	"strings"

	"github.com/kward/golib/math"
	kstrings "github.com/kward/golib/strings"
)

const MAX_COLS = 100

type tableOptions struct {
	commentPrefix  string
	ignoreComments bool
}

// CommentPrefix is an option for NewTable() that sets the comment prefix.
func CommentPrefix(v string) func(*tableOptions) error {
	return func(o *tableOptions) error { return o.setCommentPrefix(v) }
}

func (o *tableOptions) setCommentPrefix(v string) error {
	o.commentPrefix = v
	return nil
}

// IgnoreComments is a NewTable() that configures comment ignoring.
func IgnoreComments(v bool) func(*tableOptions) error {
	return func(o *tableOptions) error { return o.setIgnoreComments(v) }
}

func (o *tableOptions) setIgnoreComments(v bool) error {
	o.ignoreComments = v
	return nil
}

type Table struct {
	opts     *tableOptions
	records  [][]string
	colCount int   // Total number of columns.
	colSizes []int // The maximum size of each column.
}

func NewTable(opts ...func(*tableOptions) error) (*Table, error) {
	o := &tableOptions{}
	o.setCommentPrefix("#")
	o.setIgnoreComments(true)
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	return &Table{opts: o}, nil
}

func (t *Table) splitRow(row string, ifs string, columns int) []string {
	if t.opts.ignoreComments && strings.HasPrefix(row, t.opts.commentPrefix) {
		return []string{row}
	}

	if columns == 0 {
		columns = -1
	}
	return kstrings.SplitNMerged(row, ifs, columns)
}

func (t *Table) Split(records []string, ifs string, columns int) {
	var (
		cols    []string
		rows    [][]string
		numCols int
	)

	colSizes := make([]int, 0, MAX_COLS)

	for _, row := range records {
		cols = t.splitRow(row, ifs, columns)
		if len(cols) == 1 && strings.HasPrefix(cols[0], t.opts.commentPrefix) {
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
				colSizes[i] = math.Max(len(col), colSizes[i])
			}
		}
		rows = append(rows, cols)
	}

	t.records = rows
	t.colCount = len(colSizes)
	t.colSizes = colSizes
}

func (t *Table) IsComment(row []string) bool {
	return len(row) == 1 && strings.HasPrefix(row[0], t.opts.commentPrefix)
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
