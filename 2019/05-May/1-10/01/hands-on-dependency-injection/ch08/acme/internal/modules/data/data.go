package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gy-kim/golang-daily-practice/2019/05-May/1-10/01/hands-on-dependency-injection/ch08/acme/internal/logging"
)

const (
	// default person id (returned on error)
	defaultPersonID = 0

	// SQL statements as constants (to reduce duplication and maintenance in tests)
	sqlAllColumns = "id, fullname, phone, currency, price"
	sqlInsert     = "INSERT INTO person (fullname, phone, currrency, price) VALUES (?, ?, ?, ?)"
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
	// Logger returns a reference to the logger
	Logger() logging.Logger

	// DataDSN returns the data source name
	DataDSN() string
}

var getDB = func(cfg Config) (*sql.DB, error) {
	if db == nil {
		var err error
		db, err = sql.Open("mysql", cfg.DataDSN())
		if err != nil {
			// if the DB cannot be accessed we are dead
			panic(err.Error())
		}
	}
	return db, nil
}

// Person is the data transfer object (DTO) for this package
type Person struct {
	// ID is the unique ID for this person
	ID int

	// FullName is the name of this person
	FullName string

	// Phone is the phone for this person
	Phone string

	// Currency is the currency this person has paid in
	Currency string

	// Price is the amount (in the above currency) paid by this person
	Price float64
}

// Save will save the supplied person and return the ID of the newly created person or an error.
// Errors returned are caused by the underlying database or our connection to it.
func Save(ctx context.Context, cfg Config, in *Person) (int, error) {
	db, err := getDB(cfg)
	if err != nil {
		cfg.Logger().Error("failed to get DB connection. err: %s", err)
		return defaultPersonID, err
	}

	// set latency budget for the database call
	subCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	result, err := db.ExecContext(subCtx, sqlInsert, in.FullName, in.Phone, in.Currency, in.Price)
	if err != nil {
		cfg.Logger().Error("failed to save person into DB. err: %s", err)
		return defaultPersonID, err
	}

	// retrieve and return the ID of the person created
	id, err := result.LastInsertId()
	if err != nil {
		cfg.Logger().Error("faield to retrieve id of last saved pserson. err: %s", err)
		return defaultPersonID, err
	}

	return int(id), nil
}

// func LoadAll(ctx context.Context, cfg Config) ([]*Person, error) {

// }
