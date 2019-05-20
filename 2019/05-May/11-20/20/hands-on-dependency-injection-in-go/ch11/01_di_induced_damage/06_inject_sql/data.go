package data

import (
	"database/sql"
	"errors"
)

const (
	// default person id (returned on error)
	defaultPersonID = 0

	sqlAllColumns = "id, fullname, phone, currency, price"
	sqlInsert     = "INSERT INTO person (fullname, phone, currency, price) VALUES (?, ?, ?, ?)"
	sqlLoadAll    = "SELECT " + sqlAllColumns + " FROM person"
	sqlLoadByID   = "SELECT " + sqlAllColumns + " FROM person WHERE id = ? LIMIT 1"
)

var (
	db *sql.DB

	// ErrNotFound is returned when the no records where matched by the query
	ErrNotFound = errors.New("not found")
)

// Config is the configuration for the data package
type Config interface {
	// DataDSN returns the data source name
	DataDSN() string
}

var getDB = func(cfg Config) (Database, error) {
	if db == nil {
		var err error
		db, err = sql.Open("mysql", cfg.DataDSN())
		if err != nil {
			// if the DB cannot be accessed we are dead
			panic(err.Error())
		}
	}

	return &DatabaseImpl{db: db}, nil
}

// Person is the data transfer object (DTO) for this package
type Person struct {
	ID       int
	FullName string
	Phone    string
	Currency string
	Price    float64
}

type scanner func(desc ...interface{}) error

func populatePerson(scanner scanner) (*Person, error) {
	out := &Person{}
	err := scanner(&out.ID, &out.FullName, &out.Phone, &out.Currency, &out.Price)
	return out, err
}
