package data

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/21-31/26-31/hands-on-dependency-injection-in-go/ch09/acme/internal/logging"
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

// custom type so we can convert sql result to easily
type scanner func(desc ...interface{}) error

func populatePerson(scanner scanner) (*Person, error) {
	out := &Person{}
	err := scanner(&out.ID, &out.FullName, &out.Phone, &out.Currency, &out.Price)
	return out, err
}
