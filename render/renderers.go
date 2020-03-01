package render

import (
	"bytes"
	"encoding/csv"
	"strings"

	kstrings "github.com/kward/golib/strings"
	"github.com/kward/tabulate/table"
)

// Renderers holds a populated list of renderers.
var Renderers = []Renderer{
	&CSVRenderer{},
	&MarkdownRenderer{},
	&MySQLRenderer{},
	&PlainRenderer{},
	&SQLite3Renderer{},
}

// Renderer is an interface that allows the contents of a Table to be rendered.
type Renderer interface {
	// Render the table.
	Render(*table.Table) string
	// Type returns the type of renderer.
	Type() string

	// SectionsSupported returns true if sections are supported.
	SectionsSupported() bool
}

// MySQLRenderer implements table rendering as CSV.
type CSVRenderer struct{}

// Ensure the Renderer interface is implemented.
var _ Renderer = new(CSVRenderer)

// Render implements the Renderer interface.
func (r *CSVRenderer) Render(tbl *table.Table) string {
	if tbl == nil || tbl.NumRows() == 0 {
		return ""
	}

	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)

	for _, row := range tbl.Rows() {
		if row.IsComment() {
			continue
		}
		w.Write(row.Values())
	}
	w.Flush()
	return buf.String()
}

// Type implements the Renderer interface.
func (r *CSVRenderer) Type() string { return "csv" }

// SectionsSupported implements the Renderer interface.
func (r *CSVRenderer) SectionsSupported() bool { return false }

// MarkdownRenderer implements table rendering in Markdown format.
type MarkdownRenderer struct{}

// Ensure the Renderer interface is implemented.
var _ Renderer = new(MarkdownRenderer)

// Render implements the Renderer interface.
func (r *MarkdownRenderer) Render(tbl *table.Table) string {
	if tbl == nil || tbl.NumRows() == 0 {
		return ""
	}

	s := ""
	for _, row := range tbl.Rows() {
		if row.IsComment() {
			continue
		}

		for j, col := range row.Columns() {
			if j == 0 {
				s += "|"
			}
			if j < row.NumColumns() {
				s += " "
				s += kstrings.Stretch(col.Value(), ' ', row.Sizes()[j])
				s += " |"
			}
		}
		s += "\n"
	}
	return s
}

// Type implements the Renderer interface.
func (r *MarkdownRenderer) Type() string { return "markdown" }

// SectionsSupported implements the Renderer interface.
func (r *MarkdownRenderer) SectionsSupported() bool { return false }

// MySQLRenderer implements table rendering similar to MySQL.
type MySQLRenderer struct{}

// Ensure the Renderer interface is implemented.
var _ Renderer = new(MySQLRenderer)

// Render implements the Renderer interface.
func (r *MySQLRenderer) Render(tbl *table.Table) string {
	if tbl == nil || tbl.NumRows() == 0 {
		return ""
	}

	s := ""

	sectionBreak := "+"
	for _, size := range tbl.Rows()[0].Sizes() {
		sectionBreak += kstrings.Stretch("", '-', size+2)
		sectionBreak += "+"
	}
	sectionBreak += "\n"

	for _, row := range tbl.Rows() {
		if row.IsComment() {
			continue
		}

		for j, col := range row.Columns() {
			if j == 0 {
				s += "|"
			}
			if j < row.NumColumns() {
				s += " "
				s += kstrings.Stretch(col.Value(), ' ', row.Sizes()[j])
				s += " |"
			}
		}
		s += "\n"
	}

	if s != "" {
		s = sectionBreak + s + sectionBreak
	}

	return s
}

// Type implements the Renderer interface.
func (r *MySQLRenderer) Type() string { return "mysql" }

// SectionsSupported implements the Renderer interface.
func (r *MySQLRenderer) SectionsSupported() bool { return false }

// PlainRenderer implements table rendering as rows and columns of text.
type PlainRenderer struct {
	ofs string
}

// Ensure the Renderer interface is implemented.
var _ Renderer = new(PlainRenderer)

// Render implements the Renderer interface.
func (r *PlainRenderer) Render(tbl *table.Table) string {
	if tbl == nil || tbl.NumRows() == 0 {
		return ""
	}

	var buf bytes.Buffer
	for _, row := range tbl.Rows() {
		if row.IsComment() {
			buf.WriteString(row.Columns()[0].Value() + "\n")
			continue
		}

		tail := "" // Tail to append on *next* loop.
		for j, col := range row.Columns() {
			if col.Length() == 0 { // If this col is empty, remaining cols will be too.
				break
			}
			if j > 0 {
				tail += r.ofs
			}
			buf.WriteString(tail + col.Value())
			if j < row.NumColumns()-1 {
				tail = strings.Repeat(" ", tbl.ColSizes()[j]-col.Length())
			}
		}
		buf.WriteRune('\n')
	}
	return buf.String()
}

// Type implements the Renderer interface.
func (r *PlainRenderer) Type() string { return "plain" }

// SectionsSupported implements the Renderer interface.
func (r *PlainRenderer) SectionsSupported() bool { return true }

// SetOFS sets the OFS separator.
func (r *PlainRenderer) SetOFS(ofs string) { r.ofs = ofs }

// MySQLRenderer implements table rendering similar to SQLite3.
type SQLite3Renderer struct{}

// Ensure the Renderer interface is implemented.
var _ Renderer = new(SQLite3Renderer)

// Render implements the Renderer interface.
func (r *SQLite3Renderer) Render(tbl *table.Table) string {
	if tbl == nil || tbl.NumRows() == 0 {
		return ""
	}

	s := ""
	for _, row := range tbl.Rows() {
		if row.IsComment() {
			// Do nothing.
			continue
		}

		for j, col := range row.Columns() {
			if j > 0 {
				s += "|"
			}
			s += col.Value()
		}
		s += "\n"
	}
	return s
}

// Type implements the Renderer interface.
func (r *SQLite3Renderer) Type() string { return "sqlite3" }

// SectionsSupported implements the Renderer interface.
func (r *SQLite3Renderer) SectionsSupported() bool { return false }
