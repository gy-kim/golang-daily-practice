package data

import "context"

// NewDAO will initialize the database connection pool (if not already done) and return a data access object which
// can be used to interest with the database
func NewDAO(cfg Config) *DAO {
	// initialize the db connection pool
	_, _ = getDB(cfg)

	return &DAO{
		cfg: cfg,
	}
}

// DAO is a data access object that provides an abstraction over our database interactions.
type DAO struct {
	cfg Config

	// Tracker is an optional query timer
	Tracker QueryTracker
}

// Load will attempt to load and
func (d *DAO) Load(ctx context.Context, ID int) (*Person, error) {
	// track processing time
	return nil, nil
}
