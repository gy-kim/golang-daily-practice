package sqlmock

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
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
	queryBasedExpectation
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

func (e *ExpectedQuery) String() string {
	msg := "ExpectedQuery => expecting Query, QueryContext or QueryRow which:"
	msg += "\n  - matches sql: '" + e.expectSQL + "'"

	if len(e.args) == 0 {
		msg += "\n - is withoiut arguments"
	} else {
		msg += "\n - is with arguments:\n"
		for i, arg := range e.args {
			msg += fmt.Sprintf("    %d - %+v\n", i, arg)
		}
		msg = strings.TrimSpace(msg)
	}

	if e.rows != nil {
		msg += fmt.Sprintf("\n  = %s", e.rows)
	}

	if e.err != nil {
		msg += fmt.Sprintf("\n - should return error: %s", e.err)
	}
	return msg
}

// WillReturnRows specifies the set of resulting rows that will be returned
func (e *ExpectedQuery) WillReturnRows(rows ...*Rows) *ExpectedQuery {
	sets := make([]*Rows, len(rows))
	for i, r := range rows {
		sets[i] = r
	}
	e.rows = &rowSets{sets: sets, ex: e}
	return e
}

// ExpectedExec is used to manage *sql.DB.Exec, *sql.Tx.Exec or *sql.Stmt.Exec expectations.
// Returned by *Sqlmock.ExpectExec.
type ExpectedExec struct {
	queryBasedExpectation
	result driver.Result
	delay  time.Duration
}

// WithArgs will match given expcted args to actual database exec operation arguments.
// if at least one argument does not match, it will return an error. For specific
// arguments an sqlmock.Argument interface can be used to match an argument.
func (e *ExpectedExec) WithArgs(args ...driver.Value) *ExpectedExec {
	e.args = args
	return e
}

// WillReturnError allows to set an error for expected database exec action
func (e *ExpectedExec) WillReturnError(err error) *ExpectedExec {
	e.err = err
	return e
}

// WillDelayFor allows to specify duration for which it will delay result.
// May be used together with context.
func (e *ExpectedExec) WillDelayFor(duration time.Duration) *ExpectedExec {
	e.delay = duration
	return e
}

// String returns string representation
func (e *ExpectedExec) String() string {
	msg := "ExepectedExec => expecting Exec or ExecContext which:"
	msg += "\n  - matches sql: '" + e.expectSQL + "'"

	if len(e.args) == 0 {
		msg += "\n  - is without arguments"
	} else {
		msg += "\n  - is with arguments:\n"
		var margs []string
		for i, arg := range e.args {
			margs = append(margs, fmt.Sprintf("		%d - %+v", i, arg))
		}
		msg += strings.Join(margs, "\n")
	}

	if e.result != nil {
		res, _ := e.result.(*result)
		msg += "\n - should return Result having:"
		msg += fmt.Sprintf("\n	   LastInsertId: %d", res.insertID)
		msg += fmt.Sprintf("\n     RowsAffected: %d", res.rowsAffected)
		if res.err != nil {
			msg += fmt.Sprintf("\n    Error: %s", res.err)
		}
	}

	if e.err != nil {
		msg += fmt.Sprintf("\n  - should return error: %s", e.err)
	}
	return msg
}

// WillReturnResult arranges for an expected Exec() to return a particular
// result, there is sqlmock.NewResult(lastInsertID int64, affectedRows int64) method
// to build a corresponding result.
func (e *ExpectedExec) WillReturnResult(result driver.Result) *ExpectedExec {
	e.result = result
	return e
}

type ExpectedPrepare struct {
	commonExpectation
	mock         *sqlmock
	expectSQL    string
	statement    driver.Stmt
	closeErr     error
	mustBeClosed bool
	wasClosed    bool
	delay        time.Duration
}

// WillReturnError allows to set an error for the expected *sql.Prepare or *sql.Tx.Prepare action.
func (e *ExpectedPrepare) WillReturnError(err error) *ExpectedPrepare {
	e.err = err
	return e
}

// WillReturnCloseError allows to set an error for this prepared statement Close action.
func (e *ExpectedPrepare) WillReturnCloseError(err error) *ExpectedPrepare {
	e.closeErr = err
	return e
}

// WillDelayFor allows to specify duration for which it will delay result.
// May be used together with Context
func (e *ExpectedPrepare) WillDelayFor(duration time.Duration) *ExpectedPrepare {
	e.delay = duration
	return e
}

// WillBeClosed expects this prepared statement to be closed.
func (e *ExpectedPrepare) WillBeClosed() *ExpectedPrepare {
	e.mustBeClosed = true
	return e
}

// ExpectQuery allows to expect Query() or QueryRow() on this prepared statement.
func (e *ExpectedPrepare) ExpectQuery() *ExpectedQuery {
	eq := &ExpectedQuery{}
	eq.expectSQL = e.expectSQL
	eq.converter = e.mock.converter
	return eq
}

// ExpectExec allows to expect Exec() on this prepared statement.
func (e *ExpectedPrepare) ExpectExec() *ExpectedExec {
	eq := &ExpectedExec{}
	eq.expectSQL = e.expectSQL
	eq.converter = e.mock.converter
	e.mock.expected = append(e.mock.expected, eq)
	return eq
}

// String returns string representation
func (e *ExpectedPrepare) String() string {
	msg := "ExpectedPrepare => expecting Prepare statement which:"
	msg += "\n - matches sql: '" + e.expectSQL + "'"

	if e.err != nil {
		msg += fmt.Sprintf("\n - should return error: %s", e.err)
	}

	if e.closeErr != nil {
		msg += fmt.Sprintf("\n - should return error on Close: %s", e.closeErr)
	}

	return msg
}

// query based expectation
// adds a query matching logic
type queryBasedExpectation struct {
	commonExpectation
	expectSQL string
	converter driver.ValueConverter
	args      []driver.Value
}

func (e *queryBasedExpectation) attemptArgMatch(args []namedValue) (err error) {
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

func (e *queryBasedExpectation) argsMatches(args []namedValue) error {
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
