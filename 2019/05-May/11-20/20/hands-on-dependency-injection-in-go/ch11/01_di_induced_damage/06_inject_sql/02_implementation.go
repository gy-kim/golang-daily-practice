package data

import (
	"context"
	"database/sql"
)

type DatabaseImpl struct {
	db *sql.DB
}

func (c *DatabaseImpl) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	return c.db.QueryRowContext(ctx, query, args)
}

func (c *DatabaseImpl) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	return c.db.QueryContext(ctx, query, args...)
}

func (c *DatabaseImpl) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

type RowImpl struct {
	rows *sql.Rows
}

func (r RowImpl) Scan(desc ...interface{}) error {
	return r.rows.Scan(desc)
}

type RowsImpl struct {
	rows *sql.Rows
}

func (r RowsImpl) Scan(desc ...interface{}) error {
	return r.rows.Scan(desc...)
}

func (r RowsImpl) Close() error {
	return r.rows.Close()
}

func (r RowsImpl) Next() bool {
	return r.rows.Next()
}

type ResultImpl struct {
	result sql.Result
}

func (r *ResultImpl) LastInsertId() (int64, error) {
	return r.result.LastInsertId()
}

func (r *ResultImpl) RowsAffected() (int64, error) {
	return r.result.RowsAffected()
}
