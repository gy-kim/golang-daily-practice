package data

import (
	"database/sql"
	"errors"

	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/01-10/08-09/hands-on-dependency-injection-in-go/ch04/acme/internal/config"
	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/01-10/08-09/hands-on-dependency-injection-in-go/ch04/acme/internal/logging"
)

const (
	defaultPersonID = 0
)

var (
	db *sql.DB

	ErrNotFound = errors.New("not found")
)

func getDB() (*sql.DB, error) {
	if db == nil {
		if config.App == nil {
			return nil, errors.New("config is not initialized")
		}

		var err error
		db, err = sql.Open("mysql", config.App.DSN)
		if err != nil {
			panic(err.Error())
		}
	}

	return db, nil
}

type Person struct {
	ID       int
	FullName string
	Phone    string
	Currency string
	Price    float64
}

func Save(in *Person) (int, error) {
	db, err := getDB()
	if err != nil {
		logging.L.Error("failed to get DB connection. err: %s", err)
		return defaultPersonID, err
	}

	// perform DB insert
	query := "INSERT INTO person (fullname, phone, currrency, price) VALUES(?, ?, ?, ?)"
	result, err := db.Exec(query, in.FullName, in.Phone, in.Currency, in.Price)
	if err != nil {
		logging.L.Error("failed to save person into DB. err: %s", err)
		return defaultPersonID, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logging.L.Error("failed to retrieve id of last saved person. err: %s", err)
		return defaultPersonID, err
	}
	return int(id), nil
}

func LoadAll() ([]*Person, error) {
	db, err := getDB()
	if err != nil {
		logging.L.Error("failed to get DB connection. err: %s", err)
		return nil, err
	}

	// perform DB select
	query := "SELECT id, fullname, phone, currency, FROM person"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var out []*Person

	for rows.Next() {
		record, err := populatePerson(rows.Scan)
		if err != nil {
			logging.L.Error("failed to convert query result. err: %s", err)
			return nil, err
		}

		out = append(out, record)
	}

	if len(out) == 0 {
		logging.L.Warn("no people foubnd in the database")
		return nil, ErrNotFound
	}

	return out, nil
}

func Load(ID int) (*Person, error) {
	db, err := getDB()
	if err != nil {
		logging.L.Error("failed to get DB connection. err: %s", err)
		return nil, err
	}

	// perform DB select
	query := "SELECT id, fullname , phone, price FROM person WHERE id = ? LIMIT 1"
	row := db.QueryRow(query, ID)

	// retrieve columns and populate the person object
	out, err := populatePerson(row.Scan)
	if err != nil {
		if err == sql.ErrNoRows {
			logging.L.Warn("failed to load requested person '%d'. err: %s", ID, err)
			return nil, ErrNotFound
		}

		logging.L.Error("failed to convert query result. err: %s", err)
		return nil, err
	}
	return out, nil
}

type scanner func(dest ...interface{}) error

func populatePerson(scanner scanner) (*Person, error) {
	out := &Person{}
	err := scanner(&out.ID, &out.FullName, &out.Phone, &out.Currency, &out.Price)
	return out, err
}
