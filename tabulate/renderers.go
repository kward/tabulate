package tabulate

import (
	"bytes"
	"encoding/csv"
	"strings"

	kstrings "github.com/kward/golib/strings"
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
	Render(*Table) string
	// Type returns the type of renderer.
	Type() string
}

type CSVRenderer struct{}

func (r *CSVRenderer) Render(tbl *Table) string {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)

	for _, record := range tbl.records {
		w.Write(record)
	}
	w.Flush()
	return buf.String()
}

func (r *CSVRenderer) Type() string {
	return "csv"
}

type MarkdownRenderer struct{}

func (r *MarkdownRenderer) Render(tbl *Table) string {
	s := ""

	for _, row := range tbl.records {
		if tbl.IsComment(row) {
			continue
		}

		for colNum, col := range row {
			if colNum == 0 {
				s += "|"
			}
			if colNum < tbl.colCount {
				s += " "
				s += kstrings.Stretch(col, ' ', tbl.colSizes[colNum])
				s += " |"
			}
		}
		s += "\n"
	}

	return s
}

func (r *MarkdownRenderer) Type() string {
	return "markdown"
}

type MySQLRenderer struct{}

func (r *MySQLRenderer) Render(tbl *Table) string {
	s := ""

	sectionBreak := "+"
	for _, colSize := range tbl.colSizes {
		sectionBreak += kstrings.Stretch("", '-', colSize+2)
		sectionBreak += "+"
	}
	sectionBreak += "\n"

	for _, row := range tbl.records {
		if tbl.IsComment(row) {
			continue
		}

		for colNum, col := range row {
			if colNum == 0 {
				s += "|"
			}
			if colNum < tbl.colCount {
				s += " "
				s += kstrings.Stretch(col, ' ', tbl.colSizes[colNum])
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

func (r *MySQLRenderer) Type() string {
	return "mysql"
}

type PlainRenderer struct {
	OFS string
}

func (r *PlainRenderer) Render(tbl *Table) string {
	s := ""

	for _, row := range tbl.records {
		if tbl.IsComment(row) {
			s += row[0] + "\n"
			continue
		}

		tail := "" // Tail to append on *next* loop.
		for colNum, col := range row {
			if len(col) == 0 { // If this col is empty, remaining cols will be too.
				break
			}
			if colNum > 0 {
				tail += r.OFS
			}
			s += tail + col
			if colNum < tbl.colCount-1 {
				tail = strings.Repeat(" ", tbl.colSizes[colNum]-len(col))
			}
		}
		s += "\n"
	}
	return s
}

func (r *PlainRenderer) Type() string {
	return "plain"
}

type SQLite3Renderer struct{}

func (r *SQLite3Renderer) Render(tbl *Table) string {
	s := ""

	for _, row := range tbl.records {
		if tbl.IsComment(row) {
			// Do nothing.
			continue
		}

		for colNum, col := range row {
			if colNum > 0 {
				s += "|"
			}
			s += col
		}
		s += "\n"
	}
	return s
}

func (r *SQLite3Renderer) Type() string {
	return "sqlite3"
}
