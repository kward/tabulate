/*
Package table provides functionality for holding and describing tabular data.

TODO(2020-03-01) Add support for justification.
- The table should have a default justification for each column.
- Each column should be able to override the default table justification.
*/
package table

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kward/golib/math"
	kstrings "github.com/kward/golib/strings"
)

const MAX_COLS = 100

// Row describes a row in the table.
type Row struct {
	columns   []*Column // Columnar data of the row.
	sizes     []int     // Sizes of the columns.
	isComment bool
}

// NewRow instantiates a new row. If the row is a comment, there can be only one
// record.
func NewRow(records []string, isComment bool) (*Row, error) {
	if isComment && len(records) > 1 {
		return nil, fmt.Errorf("only one record allowed for a comment row")
	}
	return newRow(records, isComment), nil
}

func newRow(records []string, isComment bool) *Row {
	cols := []*Column{}
	sizes := []int{}
	for _, r := range records {
		cols = append(cols, &Column{cell: r})
		sizes = append(sizes, len(r))
	}
	return &Row{cols, sizes, isComment}
}

// Values returns the cell data for the row.
func (r *Row) Values() []string {
	vs := make([]string, r.NumColumns())
	for i, c := range r.Columns() {
		vs[i] = c.Value()
	}
	return vs
}

// Columns held in the row.
func (r *Row) Columns() []*Column { return r.columns }

// NumColumns returns the number of columns in the row.
func (r *Row) NumColumns() int { return len(r.columns) }

// Sizes of the columns.
func (r *Row) Sizes() []int { return r.sizes }

// IsComment returns true if the full line is a comment.
func (r *Row) IsComment() bool { return r.isComment }

// String implements fmt.Stringer.
func (r *Row) String() string {
	var buf bytes.Buffer
	for i, col := range r.columns {
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%q", col.String()))
	}
	return buf.String()
}

// Column holds the data for each column.
type Column struct {
	cell string // The actual cell data.
}

// Value of the column.
func (c *Column) Value() string { return c.cell }

// Length of the cell.
func (c *Column) Length() int { return len(c.cell) }

// String implements fmt.Stringer.
func (c *Column) String() string { return c.cell }

type Table struct {
	opts *options

	rows     []*Row
	colSizes []int
}

func NewTable(opts ...func(*options) error) (*Table, error) {
	o := &options{}
	o.setCommentPrefix("#")
	o.setEnableComments(false)
	o.setSectionReset(false)
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	return &Table{
		opts: o,
		rows: []*Row{},
	}, nil
}

// Append lines to the table.
func (t *Table) Append(records ...[]string) {
	for _, rs := range records {
		cs := make([]*Column, len(rs)) // Columns.
		for i, r := range rs {
			cs[i] = &Column{cell: r}
		}
		t.rows = append(t.rows, &Row{
			columns: cs,
		})
	}
}

// ColSizes returns the maximum size of each column.
func (t *Table) ColSizes() []int { return t.colSizes }

// Rows returns the table row data.
func (t *Table) Rows() []*Row { return t.rows }

// NumRows returns the number of rows in the table.
func (t *Table) NumRows() int { return len(t.rows) }

// String implements fmt.Stringer.
func (t *Table) String() string {
	var buf bytes.Buffer
	buf.WriteRune('[')
	for i, row := range t.rows {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(row.String())
	}
	buf.WriteRune(']')
	return fmt.Sprintf("%v sizes: %v\n", buf.String(), t.colSizes)
}

// Split lines of text into a table.
//
// The count `n` determines the number of substrings to return:
//
//     n > 0: at most n columns; the last column will be the unsplit remainder.
//     n == 0: the result is nil (an empty table)
//     n < 0: all columns
func Split(lines []string, ifs string, n int, opts ...func(*options) error) (*Table, error) {
	if n > MAX_COLS {
		return nil, fmt.Errorf("column count %d exceeds supported maximum number of columns %d", n, MAX_COLS)
	}

	tbl, err := NewTable(opts...)
	if err != nil {
		return nil, fmt.Errorf("error instantiating a table; %s", err)
	}

	var sizes []int
	switch {
	case n > 0:
		sizes = make([]int, n)
	case n == 0:
		return tbl, nil
	case n < 0:
		sizes = make([]int, MAX_COLS)
	}

	colsSeen := 0
	rows := []*Row{}
	for _, line := range lines {
		row := splitLine(tbl.opts, line, ifs, n)
		for j, col := range row.Columns() {
			if j >= MAX_COLS-1 {
				j = MAX_COLS - 1
			}
			if !row.IsComment() {
				sizes[j] = math.Max(col.Length(), sizes[j])
			}
		}
		rows = append(rows, row)
		colsSeen = math.Max(row.NumColumns(), colsSeen)
	}

	tbl.rows = rows
	tbl.colSizes = sizes[:colsSeen]
	return tbl, nil
}

func splitLine(opts *options, line string, ifs string, columns int) *Row {
	isComment := false
	var recs []string
	if opts.enableComments && strings.HasPrefix(line, opts.commentPrefix) {
		recs = []string{line}
		isComment = true
	} else {
		recs = kstrings.SplitNMerged(line, ifs, columns)
	}
	return newRow(recs, isComment)
}
