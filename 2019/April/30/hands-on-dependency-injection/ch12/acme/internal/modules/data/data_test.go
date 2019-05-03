package data

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/gy-kim/golang-daily-practice/2019/April/30/hands-on-dependency-injection/ch12/acme/internal/logging"

	"github.com/stretchr/testify/require"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSve_happyPath(t *testing.T) {
	// define context and therefore test timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// define a mock db
	testDb, dbMock, err := sqlmock.New()
	defer testDb.Close()
	require.NoError(t, err)

	// configure the mock db
	queryRegex := convertSQLToRegex(sqlInsert)
	dbMock.ExpectExec(queryRegex).WillReturnResult(sqlmock.NewResult(2, 1))

	// monkey patching strts here
	db = testDb
	// end of money patch

	// inputs
	in := &Person{
		FullName: "Jake Blues",
		Phone:    "0123456789",
		Currency: "AUD",
		Price:    123.45,
	}

	// call function

}

func TestLoad_tableDrivenTest(t *testing.T) {
	// define context and therefore test timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	scenarios := []struct {
		desc            string
		configureMockDB func(sqlmock.Sqlmock)
		expectedResult  *Person
		expectError     bool
	}{
		{
			desc: "happy path",
			configureMockDB: func(dbMock sqlmock.Sqlmock) {
				queryRegex := convertSQLToRegex(sqlLoadAll)
				dbMock.ExpectQuery(queryRegex).WillReturnRows(
					sqlmock.NewRows(strings.Split(sqlAllColumns, ", ")).
						AddRow(2, "Paul", "0123456789", "CAD", 23.45))
			},
			expectedResult: &Person{
				ID:       2,
				FullName: "Paul",
				Phone:    "0123456789",
				Currency: "CAD",
				Price:    23.45,
			},
			expectError: false,
		}, {
			desc: "load error",
			configureMockDB: func(dbMock sqlmock.Sqlmock) {
				queryRegex := convertSQLToRegex(sqlLoadAll)
				dbMock.ExpectQuery(queryRegex).WillReturnError(errors.New("something failed"))
			},
			expectedResult: nil,
			expectError:    true,
		},
	}

	for _, scenario := range scenarios {
		// define a mock db
		testDb, dbMock, err := sqlmock.New()
		require.NoError(t, err)

		// configure the mock db
		scenario.configureMockDB(dbMock)

		// monkey db for this test
		original := *db
		db = testDb

		// call function
	}
}

func convertSQLToRegex(in string) string {
	return `\Q` + in + `\E`
}

type testConfig struct{}

// Logger implements Config
func (t *testConfig) Logger() logging.Logger {
	return logging.LoggerStdOut{}
}

func (t *testConfig) DataDSN() string {
	return ""
}
