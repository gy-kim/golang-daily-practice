package sqlmock

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// an expectation interface
type expectation interface {
	fulfilled() bool
	Lock()
	Unlock()
	String() string
}

// common expectation struct
// satisfies the expectation interface
type commonExpectation struct {
	sync.Mutex
	triggered bool
	err       error
}

func (e *commonExpectation) fulfilled() bool {
	return e.triggered
}

// ExpectedClose is used to manage *sql.DB.Close expectation
// returned by *Sqlmock.ExpectedClose.
type ExpectedClose struct {
	commonExpectation
}

// WillReturnError allows to set an error for *sql.DB.Close action
func (e *ExpectedClose) WillReturnError(err error) *ExpectedClose {
	e.err = err
	return e
}

// String returns string representation
func (e *ExpectedClose) String() string {
	msg := "ExpectedClose => expecting database Close"
	if e.err != nil {
		msg += fmt.Sprintf(", which should return error: %s", e.err)
	}
	return msg
}

// ExpectedBegin is used to manage *sql.DB.Begin expectation
// returned by *Sqlmock.ExpectBegin.
type ExpectedBegin struct {
	commonExpectation
	delay time.Duration
}

// WillReturnError allows to set an error for *sql.DB.Begin action
func (e *ExpectedBegin) WillReturnError(err error) *ExpectedBegin {
	e.err = err
	return e
}

// String returns string representation.
func (e *ExpectedBegin) String() string {
	msg := "ExpectedBegin => expecting database transaction Begin"
	if e.err != nil {
		msg += fmt.Sprintf(", which should return error: %s", e.err)
	}
	return msg
}

// WillDelayFor allows to specify duration for which it will delay result.
// May be used together with Context
func (e *ExpectedBegin) WillDelayFor(duration time.Duration) *ExpectedBegin {
	e.delay = duration
	return e
}

// ExpectedCommit is used to manage *sql.Tx.Commit expectation
// returned by *Sqlmock.ExpectCommit.
type ExpectedCommit struct {
	commonExpectation
}

// WillReturnError allows to set an error for *sql.Tx.Close action
func (e *ExpectedCommit) WillReturnError(err error) *ExpectedCommit {
	e.err = err
	return e
}

// String returns string representation
func (e *ExpectedCommit) String() string {
	msg := "ExpectedCommit => expecting transaction Commit"
	if e.err != nil {
		msg += fmt.Sprintf(", which should return error: %s", e.err)
	}
	return msg
}

// ExpectedRollback is used to manage *sql.Tx.Rollback expectation
// returned by *Sqlmock.ExpectRollback.
type ExpectedRollback struct {
	commonExpectation
}

// WillReturnError allows to set an error for *sql.Tx.Rollback action
func (e *ExpectedRollback) WillReturnError(err error) *ExpectedRollback {
	e.err = err
	return e
}

// String returns string representation.
func (e *ExpectedRollback) String() string {
	msg := "ExpectedRollback => expecting transaction Rollback"
	if e.err != nil {
		msg += fmt.Sprintf(", which should return error: %s", e.err)
	}
	return msg
}

// ExpectedQuery is used to manage *sql.DB.Query, *sql.DB.QueryRow, *sql.Tx.Query,
// *sql.Tx.QueryRow, *sql.Stmt.Query or *sql.Stmt.QueryRow expectations.
type ExpectedQuery struct {
	queryBasedException
	rows             driver.Rows
	delay            time.Duration
	rowsMustBeClosed bool
	rowsWereClosed   bool
}

// WithArgs will match given expected args to actual database query arguments.
// if at least one argument does not match, it will return an error. For specific
// arguemtns an sqlmock.Argument interface can be used to match an argument.
func (e *ExpectedQuery) WithArgs(args ...driver.Value) *ExpectedQuery {
	e.args = args
	return e
}

// RowsWillBeClosed expects this query rows to be closed.
func (e *ExpectedQuery) RowsWillBeClosed() *ExpectedQuery {
	e.rowsMustBeClosed = true
	return e
}

// WillReturnError allows to set an error for expected database query
func (e *ExpectedQuery) WillReturnError(err error) *ExpectedQuery {
	e.err = err
	return e
}

// WillDelayFor allows to specify duration for which it will deplay
// result. May be used together with Context.
func (e *ExpectedQuery) WillDelayFor(duration time.Duration) *ExpectedQuery {
	e.delay = duration
	return e
}

// query based expectation
// adds a query matching logic
type queryBasedException struct {
	commonExpectation
	expectSQL string
	converter driver.ValueConverter
	args      []driver.Value
}

func (e *queryBasedException) attemptArgMatch(args []namedValue) (err error) {
	// catch panic
	defer func() {
		if e := recover(); e != nil {
			_, ok := e.(error)
			if !ok {
				err = fmt.Errorf(e.(string))
			}
		}
	}()

	err = e.argsMatches(args)
	return
}

func (e *queryBasedException) argsMatches(args []namedValue) error {
	if nil == e.args {
		return nil
	}
	if len(args) != len(e.args) {
		return fmt.Errorf("expected %d, but got %d arguments", len(e.args), len(args))
	}

	for k, v := range args {
		matcher, ok := e.args[k].(Argument)
		if ok {
			if !matcher.Match(v.Value) {
				return fmt.Errorf("matcher %T could not match %d argument %T - %+v", matcher, k, args[k], args[k])
			}
			continue
		}
		dval := e.args[k]
		if named, isNamed := dval.(sql.NamedArg); isNamed {
			dval = named.Value
			if v.Name != named.Name {
				return fmt.Errorf("named argument %d: name: \"%s\" does not match expected: \"%s\"", k, v.Name, named.Name)
			}
		} else if k+1 != v.Ordinal {
			return fmt.Errorf("argument %d: ordinal position: %d does not match expected: %d", k, k+1, v.Ordinal)
		}

		// convert to driver converter
		darg, err := e.converter.ConvertValue(dval)
		if err != nil {
			return fmt.Errorf("could not convert %d argument %T - %+v to driver value: %s", k, e.args[k], e.args[k], err)
		}

		if !reflect.DeepEqual(darg, v.Value) {
			return fmt.Errorf("argument %d expected [%T - %+v] does not match actual [%T - %+v", k, darg, darg, v.Value, v.Value)
		}
	}
	return nil
}
