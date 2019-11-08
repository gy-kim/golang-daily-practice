package sqlmock

import "database/sql/driver"

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

type namedValue struct {
	Name    string
	Ordinal int
	Value   driver.Value
}
