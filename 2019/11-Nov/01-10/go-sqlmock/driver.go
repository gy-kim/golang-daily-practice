package sqlmock

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync"
)

var pool *mockDriver

func init() {
	pool = &mockDriver{
		conns: make(map[string]*sqlmock),
	}
	sql.Register("sqlmock", pool)
}

type mockDriver struct {
	sync.Mutex
	counter int
	conns   map[string]*sqlmock
}

func (d *mockDriver) Open(dsn string) (driver.Conn, error) {
	d.Lock()
	defer d.Unlock()

	c, ok := d.conns[dsn]
	if !ok {
		return c, fmt.Errorf("expected a connection to be available, but it is not")
	}

	c.opened++
	return c, nil
}

// New creates sqlmock database connection and a mock to manage expectations.
// Accepts options, like ValueConverterOption, to use a ValueConverter from a specific driver.
func New(options ...func(*sqlmock) error) (*sql.DB, Sqlmock, error) {
	pool.Lock()
	dsn := fmt.Sprintf("sqlmock_db_%d", pool.counter)
	pool.counter++

	smock := &sqlmock{dsn: dsn, drv: pool, ordered: true}
	pool.conns[dsn] = smock
	pool.Unlock()

	return smock.open(options)
}
