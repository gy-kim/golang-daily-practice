package sqlmock

import (
	"bytes"
	"database/sql/driver"
	"encoding/csv"
	"fmt"
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

func (rs *rowSets) String() string {
	if rs.empty() {
		return "with empty rows"
	}

	msg := "should return rows:\n"
	if len(rs.sets) == 1 {
		for n, row := range rs.sets[0].rows {
			msg += fmt.Sprintf("    row %d - %+v\n", n, row)
		}
		return strings.TrimSpace(msg)
	}
	for i, set := range rs.sets {
		msg += fmt.Sprintf("	result set: %d\n", i)
		for n, row := range set.rows {
			msg += fmt.Sprintf("	row %d - %+v\n", n, row)
		}
	}
	return strings.TrimSpace(msg)
}

func (rs *rowSets) empty() bool {
	for _, set := range rs.sets {
		if len(set.rows) > 0 {
			return false
		}
	}
	return true
}

func rawBytes(col driver.Value) (_ []byte, ok bool) {
	val, ok := col.([]byte)
	if !ok || len(val) == 0 {
		return nil, false
	}

	b := make([]byte, len(val))
	copy(b, val)
	return b, true
}

func (rs *rowSets) invalidateRaw() {
	b := []byte(invalidate)
	for _, r := range rs.raw {
		copy(r, bytes.Repeat(b, len(r)/len(b)+1))
	}
	rs.raw = nil
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

// CloseError allows to set an error which will be returned by rows.Close function.
func (r *Rows) CloseError(err error) *Rows {
	r.closeErr = err
	return r
}

// RowError allows to set an error which will be returnd when a given row number is read
func (r *Rows) RowError(row int, err error) *Rows {
	r.nextErr[row] = err
	return r
}

// AddRow composed from database driver.Value slice return the same instance to perform usbsequent actions.
// Note that the number of values must match the number of columns
func (r *Rows) AddRow(values ...driver.Value) *Rows {
	if len(values) != len(r.cols) {
		panic("Expected number of values to match number of columns")
	}
	row := make([]driver.Value, len(r.cols))
	for i, v := range values {
		var err error
		v, err = r.converter.ConvertValue(v)
		if err != nil {
			panic(fmt.Errorf(
				"row #%d, column #%d (%q) type %T: %s",
				len(r.rows)+1, i, r.cols[i], values[i], err,
			))
		}

		row[i] = v
	}
	r.rows = append(r.rows, row)
	return r
}

// FromCSVString build rows from csv string.
func (r *Rows) FromCSVString(s string) *Rows {
	res := strings.NewReader(strings.TrimSpace(s))
	csvReader := csv.NewReader(res)

	for {
		res, err := csvReader.Read()
		if err != nil || res == nil {
			break
		}

		row := make([]driver.Value, len(r.cols))
		for i, v := range res {
			row[i] = CSVColumnParser(strings.TrimSpace(v))
		}
		r.rows = append(r.rows, row)
	}
	return r
}
