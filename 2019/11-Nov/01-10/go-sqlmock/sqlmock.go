package sqlmock

import (
	"database/sql"
	"database/sql/driver"
)

// Sqlmock interface serves to create expectations
type Sqlmock interface {
	// ExpectClose queues an expectation for this database
	// action to be triggered. the *ExpectedClose allows
	// to mock database response
	ExpectClose() *ExpectedClose

	// ExpectationWereMet checks whether all queued expectations
	// were met in order. If any of them was not met - an error is returned.
	ExpectationsWereMet() error
}

type sqlmock struct {
	ordered      bool
	dsn          string
	opened       int
	drv          *mockDriver
	converter    driver.ValueConverter
	queryMatcher QueryMatcher

	expected []expectation
}

func (c *sqlmock) open(options []func(*sqlmock) error) (*sql.DB, Sqlmock, error) {
	db, err := sql.Open("sqlmock", c.dsn)
	if err != nil {
		return db, c, err
	}
	for _, option := range options {
		err := option(c)
		if err != nil {
			return db, c, err
		}
	}
	if c.converter == nil {
		c.converter = driver.DefaultParameterConverter
	}
	if c.queryMatcher == nil {
		c.queryMatcher = QueryMatcherRegex
	}
	return db, c, db.Ping()
}

type namedValue struct {
	Name    string
	Ordinal int
	Value   driver.Value
}
