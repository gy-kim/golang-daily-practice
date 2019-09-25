package data

import (
	"database/sql"
	"errors"
)

const (
	// default person id (returned on error)
	defaultPersonID = 0

	// SQL statements as contants (to reduce duplication and maintenance in tests)
	sqlAllColumns = "id, fullname, phone, currency, price"
	sqlInsert     = "INSERT INTO person (fullname, phone, currency, price) VALUES (?, ?, ?, ?)"
	sqlLoadAll    = "SELECT " + sqlAllColumns + " FROM person"
	sqlLoadByID   = "SELECT " + sqlAllColumns + " FROM person where id = ? LIMIT 1"
)

var (
	db *sql.DB

	// ErrNotFound is returned when the no records where matched by the query
	ErrNotFound = errors.New("not found")
)
