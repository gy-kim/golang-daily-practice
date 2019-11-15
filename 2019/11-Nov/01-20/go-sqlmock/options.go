package sqlmock

import "database/sql/driver"

// ValueConverterOption allows to create a sqlmock connection
func ValueConverterOption(converter driver.ValueConverter) func(*sqlmock) error {
	return func(s *sqlmock) error {
		s.converter = converter
		return nil
	}
}

// QueryMatcherOption allows to customize SQL query matcher
func QueryMatcherOption(queryMatcher QueryMatcher) func(*sqlmock) error {
	return func(s *sqlmock) error {
		s.queryMatcher = queryMatcher
		return nil
	}
}
