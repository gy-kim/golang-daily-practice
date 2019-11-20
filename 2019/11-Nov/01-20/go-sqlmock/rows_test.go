package sqlmock

import (
	"bytes"
	"database/sql"
	"fmt"
	"testing"
)

const invalid = `☠☠☠ MEMORY OVERWRITTEN ☠☠☠ `

func ExampleRows() {
	db, mock, err := New()
	if err != nil {
		fmt.Println("failed to open sqlmock database:", err)
	}
	defer db.Close()

	rows := NewRows([]string{"id", "title"}).
		AddRow(1, "one").
		AddRow(2, "two")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	rs, _ := db.Query("SELECT")
	defer rs.Close()

	for rs.Next() {
		var id int
		var title string
		rs.Scan(&id, &title)
		fmt.Println("scanned id:", id, "and title:", title)
	}

	if rs.Err() != nil {
		fmt.Println("got rows error:", rs.Err())
	}
}

func ExampleRows_rowError() {
	db, mock, err := New()
	if err != nil {
		fmt.Println("failed to open sqlmock database:", err)
	}
	defer db.Close()

	rows := NewRows([]string{"id", "title"}).
		AddRow(0, "one").
		AddRow(1, "two").
		RowError(1, fmt.Errorf("row error"))
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	rs, _ := db.Query("SELECT")
	defer rs.Close()

	for rs.Next() {
		var id int
		var title string
		rs.Scan(&id, &title)
		fmt.Println("scanned id:", id, "and title:", title)
	}

	if rs.Err() != nil {
		fmt.Println("got rows error: ", rs.Err())
	}
}

func ExampleRows_closeError() {
	db, mock, err := New()
	if err != nil {
		fmt.Println("failed to open sqlmock database:", err)
	}
	defer db.Close()

	rows := NewRows([]string{"id", "title"}).CloseError(fmt.Errorf("close error"))
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	rs, _ := db.Query("SELECT")

	if err := rs.Close(); err != nil {
		fmt.Println("got error:", err)
	}
}

func ExampleRows_rawBytes() {
	db, mock, err := New()
	if err != nil {
		fmt.Println("failed to open sqlmock database:", err)
	}
	defer db.Close()

	rows := NewRows([]string{"id", "binary"}).
		AddRow(1, []byte(`one binary value with some text!`)).
		AddRow(2, []byte(`two binary value with even more text than the first one`))

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	rs, _ := db.Query("SELECT")
	defer rs.Close()

	type scanned struct {
		id  int
		raw sql.RawBytes
	}
	fmt.Println("initial read...")
	var ss []scanned
	for rs.Next() {
		var s scanned
		rs.Scan(&s.id, &s.raw)
		ss = append(ss, s)
		fmt.Println("scanned id:", s.id, "and raw:", string(s.raw))
	}

	if rs.Err() != nil {
		fmt.Println("got rows error:", rs.Err())
	}

	fmt.Println("after reading all...")
	for _, s := range ss {
		fmt.Println("scanned id:", s.id, "and raw:", string(s.raw))
	}
}

func ExampleRows_expectToBeClosed() {
	db, mock, err := New()
	if err != nil {
		fmt.Println("failed to open sqlmock datacase:", err)
	}
	defer db.Close()

	rows := NewRows([]string{"id", "title"}).AddRow(1, "john")
	mock.ExpectQuery("SELECT").WillReturnRows(rows).RowsWillBeClosed()

	db.Query("SELECT")

	if err := mock.ExpectationsWereMet(); err != nil {
		fmt.Println("got error:", err)
	}
}

func ExampleRows_customDriverValue() {
	db, mock, err := New()
	if err != nil {
		fmt.Println("failed to open sqlmock database:", err)
	}
	defer db.Close()

	rows := NewRows([]string{"id", "null_int"}).
		AddRow(1, 7).
		AddRow(5, sql.NullInt64{Int64: 5, Valid: true}).
		AddRow(2, sql.NullInt64{})

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	rs, _ := db.Query("SELECT")
	defer rs.Close()

	for rs.Next() {
		var id int
		var num sql.NullInt64
		rs.Scan(&id, &num)
		fmt.Println("scanned id:", id, "and null int64:", num)
	}

	if rs.Err() != nil {
		fmt.Println("got rows error:", rs.Err())
	}
}

func TestAllowsToSeetRowsErrors(t *testing.T) {
	t.Parallel()
	db, mock, err := New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := NewRows([]string{"id", "title"}).
		AddRow(0, "one").
		AddRow(1, "two").
		RowError(1, fmt.Errorf("error"))
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	rs, err := db.Query("SELECT")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer rs.Close()

	if !rs.Next() {
		t.Fatal("expected the first row to be available")
	}
	if rs.Err() != nil {
		t.Fatalf("unexpected error: %s", rs.Err())
	}
	if rs.Next() {
		t.Fatal("was not expecting the second row, since there should be an error")
	}
	if rs.Err() == nil {
		t.Fatal("expected an error, but got none")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRowsCloseError(t *testing.T) {
	t.Parallel()
	db, mock, err := New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := NewRows([]string{"id"}).CloseError(fmt.Errorf("close error"))
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	rs, err := db.Query("SELECT")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if err := rs.Close(); err == nil {
		t.Fatal("expected a close error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRowsClosed(t *testing.T) {
	t.Parallel()
	db, mock, err := New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("SELECT").WillReturnRows(rows).RowsWillBeClosed()

	rs, err := db.Query("SELECT")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if err := rs.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestQuerySingleRow(t *testing.T) {
	t.Parallel()
	db, mock, err := New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := NewRows([]string{"id"}).
		AddRow(1).
		AddRow(2)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	var id int
	if err := db.QueryRow("SELECT").Scan(&id); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	mock.ExpectQuery("SELECT").WillReturnRows(NewRows([]string{"id"}))
	if err := db.QueryRow("SELECT").Scan(&id); err != sql.ErrNoRows {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestQueryRowBytesInvalidatedByNext_bytesIntoRawBytes(t *testing.T) {
	t.Parallel()
	replace := []byte(invalid)
	rows := NewRows([]string{"raw"}).
		AddRow([]byte(`one binary value with some text!`)).
		AddRow([]byte(`two binary value with even more text than the first one`))
	scan := func(rs *sql.Rows) ([]byte, error) {
		var raw sql.RawBytes
		return raw, rs.Scan(&raw)
	}
	want := []struct {
		Initial  []byte
		Replaced []byte
	}{
		{Initial: []byte(`one binary value with some text!`), Replaced: replace[:len(replace)-7]},
		{Initial: []byte(`two binary value with even more text than the first one`), Replaced: bytes.Join([][]byte{replace, replace[:len(replace)-23]}, nil)},
	}
	queryRowBytesInvalidatedByNext(t, rows, scan, want)
}

func queryRowBytesInvalidatedByNext(t *testing.T, rows *Rows, scan func(*sql.Rows) ([]byte, error), want []struct {
	Initial  []byte
	Replaced []byte
}) {
	db, mock, err := New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	rs, err := db.Query("SELECT")
	if err != nil {
		t.Fatalf("failed to query rows: %s", err)
	}

	if !rs.Next() || rs.Err() != nil {
		t.Fatal("unexpected error on first row retrieval")
	}
	var count int
	for i := 0; ; i++ {
		count++
		b, err := scan(rs)
		if err != nil {
			t.Fatalf("unexpected error scanning row: %s", err)
		}
		if exp := want[i].Initial; !bytes.Equal(b, exp) {
			t.Fatalf("expected raw value to be '%s' (len:%d), but got [%T]:%s (len:%d)", exp, len(exp), b, b, len(b))
		}
		next := rs.Next()
		if exp := want[i].Replaced; !bytes.Equal(b, exp) {
			t.Fatalf("expexlcted raw value to be replaced with '%s' (len:%d) after calling Next(), but got [%T]:%s (len:%d)", exp, len(exp), b, b, len(b))
		}
		if !next {
			break
		}
	}
	if err := rs.Err(); err != nil {
		t.Fatalf("row interation failed: %s", err)
	}
	if exp := len(want); count != exp {
		t.Fatalf("incorrect number of rows exp: %d, but got %d", exp, count)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func queryRowBytesNotInvalidatedByNext(t *testing.T, rows *Rows, scan func(*sql.Rows) ([]byte, error), want [][]byte) {
	db, mock, err := New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	rs, err := db.Query("SELECT")
	if err != nil {
		t.Fatalf("failed to query rows: %s", err)
	}

	if !rs.Next() || rs.Err() != nil {
		t.Fatal("unexpected error on first row retrieval")
	}
	var count int
	for i := 0; ; i++ {
		count++
		b, err := scan(rs)
		if err != nil {
			t.Fatalf("unexpected error scanning row: %s", err)
		}
		if exp := want[i]; !bytes.Equal(b, exp) {
			t.Fatalf("expected raw value to be replaced with '%s' (len:%d) after calling Next(), but got [%T]:%s (len:%d)", exp, len(exp), b, b, len(b))
		}
		next := rs.Next()
		if exp := want[i]; !bytes.Equal(b, exp) {
			t.Fatalf("expected raw value to be replaced with '%s' (len:%d) after calling Next(), but got [%T]:%s (len:%d)", exp, len(exp), b, b, len(b))
		}
		if !next {
			break
		}
	}
	if err := rs.Err(); err != nil {
		t.Fatalf("row interation failed: %s", err)
	}
	if exp := len(want); count != exp {
		t.Fatalf("incorrect number of rows exp: %d, but got %d", exp, count)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func queryRowBytesInvalidatedByClose(t *testing.T, rows *Rows, scan func(*sql.Rows) ([]byte, error), want struct {
	Initial  []byte
	Replaced []byte
}) {
	db, mock, err := New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	rs, err := db.Query("SELECT")
	if err != nil {
		t.Fatalf("failed to query rows: %s", err)
	}

	if !rs.Next() || rs.Err() != nil {
		t.Fatal("unexpected error on first row retrieval")
	}
	b, err := scan(rs)
	if err != nil {
		t.Fatalf("unexpected error scanning row: %s", err)
	}
	if !bytes.Equal(b, want.Initial) {
		t.Fatalf("expected raw value to be '%s' (len:%d), but got [%T]:%s (len:%d)", want.Initial, len(want.Initial), b, b, len(b))
	}
	if err := rs.Close(); err != nil {
		t.Fatalf("expected raw value to be replaced with '%s' (len:%d) after calling Next(), but got [%T]:%s (len:%d)", want.Replaced, len(want.Replaced), b, b, len(b))
	}
	if err := rs.Err(); err != nil {
		t.Fatalf("row iteration failed: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func queryRowBytesNotInvalidatedByClose(t *testing.T, rows *Rows, scan func(*sql.Rows) ([]byte, error), want []byte) {
	db, mock, err := New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	rs, err := db.Query("SELECT")
	if err != nil {
		t.Fatalf("failed to query rows: %s", err)
	}

	if !rs.Next() || rs.Err() != nil {
		t.Fatal("unexpected error on first row retrieval")
	}
	b, err := scan(rs)
	if err != nil {
		t.Fatalf("unexpected error scanning row: %s", err)
	}
	if !bytes.Equal(b, want) {
		t.Fatalf("expected raw value to be '%s' (len:%d), but got [%T]:%s (len:%d)", want, len(want), b, b, len(b))
	}
	if err := rs.Close(); err != nil {
		t.Fatalf("unexpected error closing rows: %s", err)
	}
	if !bytes.Equal(b, want) {
		t.Fatalf("expected raw value to be replaced with '%s' (len:%d) after calling Next(), but got [%T]:%s (len:%d)", want, len(want), b, b, len(b))
	}
	if err := rs.Err(); err != nil {
		t.Fatalf("row interation failed: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
