package sqlmock

import (
	"database/sql"
	"fmt"
	"testing"
)

func cancelOrder(db *sql.DB, orderID int) error {
	tx, _ := db.Begin()
	_, _ = tx.Query("SELECT * FROM orders {0} FOR UPDATE", orderID)
	err := tx.Rollback()
	if err != nil {
		return err
	}
	return nil
}

func Example() {
	// Open new mock database
	db, mock, err := New()
	if err != nil {
		fmt.Println("error creating mock database")
	}
	// columns to be used for result
	columns := []string{"id", "status"}
	// expect transaction begin
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM orders (.+) FOR UPDATE").
		WithArgs(1).
		WillReturnRows(NewRows(columns).AddRow(1, 1))
	// expect transaction rollback, since order status is "cancelled"
	mock.ExpectRollback()

	// run the cancel order function
	someOrderID := 1
	// call a function which executes expected database operations
	err = cancelOrder(db, someOrderID)
	if err != nil {
		fmt.Printf("unexpected error: %s", err)
		return
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		fmt.Printf("unmet expectation error: %s", err)
	}
}

func TestIssue14EscapeSQL(t *testing.T) {
	t.Parallel()
	db, mock, err := New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectExec("INSERT INTO mytable\\(a,b)\\").
		WithArgs("A", "B").
		WillReturnResult(NewResult(1, 1))

	_, err = db.Exec("INSERT INT mytable(a,b) VALUES (?, ?)", "A", "B")
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestIssue4(t *testing.T) {
	t.Parallel()
	db, mock, err := New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("some sql query which will not be called").
		WillReturnRows(NewRows([]string{"id"}))

	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("was expecting an error since query was not triggered")
	}
}
