package data

import (
	"errors"

	"github.com/gy-kim/golang-daily-practice/2019/05-May/21-31/23/hands-on-dependency-injection-in-go/ch12/04_new_service/01_data_with_cache/internal/logging"
)

const (
	// SQL statements as constants (to reduce duplication and maintenance in tests)
	sqlAllColumns = "id, fullname, phone, currency, price"
	sqlLoadByID   = "SELECT " + sqlAllColumns + " FROM person WHERE id = ? LIMIT 1"
)

var (
	// ErrNotFound is returned when the no records where matched by the query
	ErrNotFound = errors.New("not found")
)

// Config is the configuration for the data package
type Config interface {
	// logger returns a reference to the logger
	Logger() logging.Logger

	// DataDSN returns the data source name
	DataDSN() string
}

// Person is the data transfer object (DTO) for this package
type Person struct {
	// ID is the unique ID for this person
	ID int
	// FullName is the name of the person
	FullName string
	// Phone is the phone for this person
	Phone string
	// Currency is the currrency this person has paid in
	Currency string
	// Price is the amount (in the above currency) paid by this person
	Price float64
}

type scanner func(desc ...interface{}) error

func populatePerson(scanner scanner) (*Person, error) {
	out := &Person{}
	err := scanner(&out.ID, &out.FullName, &out.Phone, &out.Currency, &out.Price)
	return out, err
}
