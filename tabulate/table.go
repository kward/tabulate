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
	sectionReset   bool
}

// CommentPrefix is an option for NewTable() that sets the comment prefix.
func CommentPrefix(v string) func(*tableOptions) error {
	return func(o *tableOptions) error { return o.setCommentPrefix(v) }
}

func (o *tableOptions) setCommentPrefix(v string) error {
	o.commentPrefix = v
	return nil
}

// IgnoreComments is a NewTable() option that configures comment ignoring.
func IgnoreComments(v bool) func(*tableOptions) error {
	return func(o *tableOptions) error { return o.setIgnoreComments(v) }
}

func (o *tableOptions) setIgnoreComments(v bool) error {
	o.ignoreComments = v
	return nil
}

// SectionReset is a NewTable() option that enables per-section column count
// resetting. Sections are delineated by empty lines.
func SectionReset(v bool) func(*tableOptions) error {
	return func(o *tableOptions) error { return o.setSectionReset(v) }
}

func (o *tableOptions) setSectionReset(v bool) error {
	o.sectionReset = v
	return nil
}

type Row struct {
	records   []string // Column data.
	sizes     []int    // Maximum size of each column.
	isComment bool
}

func NewRow(records []string) *Row {
	sizes := []int{}
	for _, rec := range records {
		sizes = append(sizes, len(rec))
	}
	return &Row{records: records, sizes: sizes}
}

// Columns returns the number of columns in the row.
func (r *Row) Columns() int { return len(r.records) }

type Table struct {
	opts *tableOptions
	rows []*Row
}

func NewTable(opts ...func(*tableOptions) error) (*Table, error) {
	o := &tableOptions{}
	o.setCommentPrefix("#")
	o.setIgnoreComments(true)
	o.setSectionReset(false)
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	return &Table{opts: o}, nil
}

// IsComment returns true if the full line is a comment.
func (t *Table) IsComment(row *Row) bool {
	if len(row.records) > 1 {
		return false
	}
	return strings.HasPrefix(row.records[0], t.opts.commentPrefix)
}

// Split `lines` by `ifs` into a maximum number of `columns`.
func (t *Table) Split(lines []string, ifs string, columns int) {
	sizes := make([]int, 0, MAX_COLS)
	if columns > 0 {
		sizes = sizes[:columns] // Grow sizes to given number of columns.
	} else {
		sizes = sizes[:MAX_COLS]
	}

	rows := []*Row{}
	for _, line := range lines {
		recs := t.splitLine(line, ifs, columns)
		row := NewRow(recs)
		row.isComment = t.IsComment(row)
		for idx, col := range row.records {
			sizes[idx] = math.Max(len(col), sizes[idx])
		}
		row.sizes = sizes
		rows = append(rows, row)
	}
	t.rows = rows
}

func (t *Table) splitLine(line string, ifs string, columns int) []string {
	if t.opts.ignoreComments && strings.HasPrefix(line, t.opts.commentPrefix) {
		return []string{line}
	}

	if columns == 0 {
		columns = -1
	}
	return kstrings.SplitNMerged(line, ifs, columns)
}
