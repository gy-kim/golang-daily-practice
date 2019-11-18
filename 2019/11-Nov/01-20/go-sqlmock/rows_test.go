package sqlmock

import "fmt"

import "database/sql"

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
