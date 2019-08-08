package data

import (
	"database/sql"
	"errors"

	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/01-10/08/hands-on-dependency-injection-in-go/ch04/acme/internal/config"
	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/01-10/08/hands-on-dependency-injection-in-go/ch04/acme/internal/logging"
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
