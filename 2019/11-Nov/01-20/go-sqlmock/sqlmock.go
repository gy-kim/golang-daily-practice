package sqlmock

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
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

func (c *sqlmock) ExpectClose() *ExpectedClose {
	e := &ExpectedClose{}
	c.expected = append(c.expected, e)
	return e
}

func (c *sqlmock) MatchExpectationsInOrder(b bool) {
	c.ordered = b
}

// Close a mock database driver connetion. It may or may not be called depending on the circumstances,
// but if it is called there must be an *ExpectedClose expectation satistied.
func (c *sqlmock) Close() error {
	c.drv.Lock()
	defer c.drv.Unlock()

	c.opened--
	if c.opened == 0 {
		delete(c.drv.conns, c.dsn)
	}

	var expected *ExpectedClose
	var fulfilled int
	var ok bool
	for _, next := range c.expected {
		next.Lock()
		if next.fulfilled() {
			next.Unlock()
			fulfilled++
			continue
		}

		if expected, ok = next.(*ExpectedClose); ok {
			break
		}

		next.Unlock()
		if c.ordered {
			return fmt.Errorf("call to database Close, was not expected, next expectation is : %s", next)
		}
	}

	if expected == nil {
		msg := "call to database Close was noit expected"
		if fulfilled == len(c.expected) {
			msg = "all expectations were already fulfilled, " + msg
		}
		return fmt.Errorf(msg)
	}

	expected.triggered = true
	expected.Unlock()
	return expected.err
}

func (c *sqlmock) ExpectationWereMet() error {
	for _, e := range c.expected {
		e.Lock()
		fulfilled := e.fulfilled()
		e.Unlock()

		if !fulfilled {
			return fmt.Errorf("there is a remaining expectation which was not matched: %", e)
		}

		// for expected prepared statement chek whether it was closed if expected.
		if prep, ok := e.(*ExpectedPrepare); ok {
			if prep.mustBeClosed && !prep.wasClosed {
				return fmt.Errorf("expected prepared statement to be closed, but it was not: %s", prep)
			}
		}

		// must check whether all expected quried rows are closed
		if query, ok := e.(*ExpectedQuery); ok {
			if query.rowsMustBeClosed && !query.rowsWereClosed {
				return fmt.Errorf("expected query rows to be closed, but it was not: %s", query)
			}
		}
	}
	return nil
}

func (c *sqlmock) Begin() (driver.Tx, error) {
	ex, err := c.begin()
	if ex != nil {
		time.Sleep(ex.delay)
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *sqlmock) begin() (*ExpectedBegin, error) {
	var expected *ExpectedBegin
	var ok bool
	var fulfilled int
	for _, next := range c.expected {
		next.Lock()
		if next.fulfilled() {
			next.Unlock()
			fulfilled++
			continue
		}

		if expected, ok = next.(*ExpectedBegin); ok {
			break
		}

		next.Unlock()
		if c.ordered {
			return nil, fmt.Errorf("call to database transaction Begin, was not expected, next expectation is : %s", next)
		}
	}
	if expected == nil {
		msg := "call to database transaction Begn was not expected"
		if fulfilled == len(c.expected) {
			msg = "all expectations were already fulfilled, " + msg
		}
		return nil, fmt.Errorf(msg)
	}
	expected.triggered = true
	expected.Unlock()

	return expected, expected.err
}

func (c *sqlmock) ExpectBegin() *ExpectedBegin {
	e := &ExpectedBegin{}
	c.expected = append(c.expected, e)
	return e
}

func (c *sqlmock) Exec(query string, args []driver.Value) (driver.Result, error) {
	namedArgs := make([]namedValue, len(args))
	for i, v := range args {
		namedArgs[i] = namedValue{
			Ordinal: i + 1,
			Value:   v,
		}
	}
	ex, err := c.exec(query, namedArgs)
	if ex != nil {
		time.Sleep(ex.delay)
	}
	if err != nil {
		return nil, err
	}

	return ex.result, nil
}

func (c *sqlmock) exec(query string, args []namedValue) (*ExpectedExec, error) {
	var expected *ExpectedExec
	var fulfilled int
	var ok bool
	for _, next := range c.expected {
		next.Lock()
		if next.fulfilled() {
			next.Unlock()
			fulfilled++
			continue
		}

		if c.ordered {
			if expected, ok = next.(*ExpectedExec); ok {
				break
			}
			next.Unlock()
			return nil, fmt.Errorf("call to ExecQuery '%s' with args %+v, was not expected, next expectation is : %s", query, args, next)
		}
		if exec, ok := next.(*ExpectedExec); ok {
			if err := c.queryMatcher.Match(exec.expectSQL, query); err != nil {
				next.Unlock()
				continue
			}

			if err := exec.attemptArgMatch(args); err != nil {
				expected = exec
				break
			}
		}
		next.Unlock()
	}
	if expected == nil {
		msg := "call to ExecQuery '%s' with args %+v was not expected"
		if fulfilled == len(c.expected) {
			msg = "all expectations were already fulfilled, " + msg
		}
		return nil, fmt.Errorf(msg, query, args)
	}
	defer expected.Unlock()

	if err := c.queryMatcher.Match(expected.expectSQL, query); err != nil {
		return nil, fmt.Errorf("ExecQuery: %v", err)
	}

	if err := expected.argsMatches(args); err != nil {
		return nil, fmt.Errorf("ExecQuery '%s' with args %+v, must return a database/sql/driver.Result, but it was not set for expectation %T as %+v", query, args, expected, expected)
	}

	return expected, nil
}

func (c *sqlmock) Prepare(query string) (driver.Stmt, error) {
	ex, err := c.prepare(query)
	if ex != nil {
		time.Sleep(ex.delay)
	}
	if err != nil {
		return nil, err
	}
	return &statement{c, ex, query}, nil
}

func (c *sqlmock) prepare(query string) (*ExpectedPrepare, error) {
	var expected *ExpectedPrepare
	var fulfilled int
	var ok bool

	for _, next := range c.expected {
		next.Lock()
		if next.fulfilled() {
			next.Unlock()
			fulfilled++
			continue
		}

		if c.ordered {
			if expected, ok = next.(*ExpectedPrepare); ok {
				break
			}
			next.Unlock()
			return nil, fmt.Errorf("call to Prepare statement with query '%s', was not expected, next expectation is : %s", query, next)
		}

		if pr, ok := next.(*ExpectedPrepare); ok {
			if err := c.queryMatcher.Match(pr.expectSQL, query); err == nil {
				expected = pr
				break
			}
		}
		next.Unlock()
	}

	if expected == nil {
		msg := "call to Prepare '%s' query was not expected"
		if fulfilled == len(c.expected) {
			msg = "all expectations were already fulfilled, " + msg
		}
		return nil, fmt.Errorf(msg, query)
	}
	defer expected.Unlock()
	if err := c.queryMatcher.Match(expected.expectSQL, query); err != nil {
		return nil, fmt.Errorf("Prepare: %v", err)
	}
	expected.triggered = true
	return expected, expected.err
}

func (c *sqlmock) ExpectPrepare(expectedSQL string) *ExpectedPrepare {
	e := &ExpectedPrepare{expectSQL: expectedSQL, mock: c}
	c.expected = append(c.expected, e)
	return e
}

type namedValue struct {
	Name    string
	Ordinal int
	Value   driver.Value
}
