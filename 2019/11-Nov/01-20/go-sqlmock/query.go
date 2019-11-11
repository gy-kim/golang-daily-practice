package sqlmock

import (
	"fmt"
	"regexp"
	"strings"
)

var re = regexp.MustCompile("\\s+")

func stripQuery(q string) (s string) {
	return strings.TrimSpace(re.ReplaceAllString(q, " "))
}

// QueryMatcher is an SQL query string matcher interface,
// which can be used to customize validation of SQL query strings.
type QueryMatcher interface {
	// Match expected SQL query string without whitespace to actual SQL.
	Match(expectedSQL, actualSQL string) error
}

// QueryMatcherFunc type is an adapter to allow the use of
// ordinary functions as QueryMatcher. If f is a function
// with the appropriate signature, QueryMatcherFunc(f) is a
// QueryMatcher that calls f.
type QueryMatcherFunc func(expectedSQL, actualSQL string) error

// Match implements the QueryMathcer
func (f QueryMatcherFunc) Match(expectedSQL, actualSQL string) error {
	return f(expectedSQL, actualSQL)
}

// QueryMatcherRegex is the default SQL querymatcher
// used by sqlmock. It parses expectedSQL to a regular
// expression and attempts to match actualSQL
var QueryMatcherRegex QueryMatcher = QueryMatcherFunc(func(expectedSQL, actualSQL string) error {
	expect := stripQuery(expectedSQL)
	actual := stripQuery(actualSQL)
	re, err := regexp.Compile(expect)
	if err != nil {
		return err
	}
	if !re.MatchString(actual) {
		return fmt.Errorf(`could not match actual sql: "%s" with expected regexp "%s"`, actual, re.String())
	}

	return nil
})

// QueryMatcherEqual is the SQL query matcher
// which simply tries a case sensitive match of
// expected and actual SQL strings without whitespace.
var QueryMatcherEqual QueryMatcher = QueryMatcherFunc(func(expectedSQL, actualSQL string) error {
	expect := stripQuery(expectedSQL)
	actual := stripQuery(actualSQL)
	if actual != expect {
		return fmt.Errorf(`actual sql: "%s" does not equal to expected "%s`, actual, expect)
	}
	return nil
})
