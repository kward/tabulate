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
