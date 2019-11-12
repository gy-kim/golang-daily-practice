package sqlmock

import (
	"database/sql/driver"
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
