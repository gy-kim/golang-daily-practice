package sqlmock

import (
	"database/sql/driver"
	"io"
	"strings"
)

const invalidate = "☠☠☠ MEMORY OVERWRITTEN ☠☠☠ "

// CSVColumnParser is a function which converts trimmed csv
// column string to a []byte representation. Currenctly transforms NULL to null
var CSVColumnParser = func(s string) []byte {
	switch {
	case strings.ToLower(s) == "null":
		return nil
	}
	return []byte(s)
}

type rowSets struct {
	sets []*Rows
	pos  int
	ex   *ExpectedQuery
	raw  [][]byte
}

func (rs *rowSets) Columns() []string {
	return rs.sets[rs.pos].cols
}

func (rs *rowSets) Close() error {
	rs.invalidateRaw()
	rs.ex.rowsWereClosed = true
	return rs.sets[rs.pos].closeErr
}

func (rs *rowSets) Next(dest []driver.Value) error {
	r := rs.sets[rs.pos]
	r.pos++
	rs.invalidateRaw()
	if r.pos > len(r.rows) {
		return io.EOF
	}

	for i, col := range r.rows[r.pos-1] {
		if b, ok := rawBytes(col); ok {
			rs.raw = append(rs.raw, b)
			dest[i] = b
			continue
		}
		dest[i] = col
	}
	return r.nextErr[r.pos-1]
}

// Rows is a mocked collection of rows to return for Query result
type Rows struct {
	converter driver.ValueConverter
	cols      []string
	rows      [][]driver.Value
	pos       int
	nextErr   map[int]error
	closeErr  error
}

// NewRows allows Rows to be created from a sql driver.Value slice or from the CSV string
// and to be used as sql driver.Rows. Use Sqlmock.NewRows instead of custom converter
func NewRows(columns []string) *Rows {
	return &Rows{
		cols:      columns,
		nextErr:   make(map[int]error),
		converter: driver.DefaultParameterConverter,
	}
}
