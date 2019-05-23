package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gy-kim/golang-daily-practice/2019/05-May/21-31/23/hands-on-dependency-injection-in-go/ch12/04_new_service/01_data_with_cache/internal/cache"
	"github.com/gy-kim/golang-daily-practice/2019/05-May/21-31/23/hands-on-dependency-injection-in-go/ch12/04_new_service/01_data_with_cache/internal/logging"
)

// DAO is a data access object that provides tan abstraction over our database interactions.
type DAO struct {
	cfg   Config
	db    *sql.DB
	cache *cache.Cache
}

// Load will attempt to load and return a person.
// It will return ErrNotFound when the requested person does not exist.
// Any other errors returned are caused by the underlying database or our connection to it.
func (d *DAO) Load(ctx context.Context, ID int) (*Person, error) {
	// load from cache
	out := d.loadFromCache(ID)
	if out != nil {
		return out, nil
	}

	// load from database
	row := d.db.QueryRowContext(ctx, sqlLoadByID, ID)

	out, err := populatePerson(row.Scan)
	// retrieve columns and populate the person object
	if err != nil {
		if err == sql.ErrNoRows {
			d.logger().Warn("failed to load requested peerson '%d'. err: %s", ID, err)
			return nil, ErrNotFound
		}

		d.logger().Error("failed to convert query result. err: %s", err)
		return nil, err
	}

	// save person in the cache
	d.saveToCache(ID, out)
	return out, nil
}

func (d *DAO) loadFromCache(ID int) *Person {
	payload, err := d.cache.Get(d.buildCacheKey(ID))
	if err != nil {
		d.logger().Error("failed to load requested person from cache with error: %s", err)
		return nil
	}

	if payload == nil {
		return nil
	}

	out := &Person{}
	err = json.Unmarshal(payload, out)
	if err != nil {
		d.logger().Error("failed to decode cache response with error: %s", err)
	}

	return out
}

func (d *DAO) saveToCache(ID int, person *Person) {
	payload, err := json.Marshal(person)
	if err != nil {
		d.logger().Error("failed to encode person to JSON with error: %s")
		return
	}

	err = d.cache.Set(d.buildCacheKey(ID), payload)
	if err != nil {
		d.logger().Error("failed to save person into cache with error: %s", err)
	}
}

func (d *DAO) buildCacheKey(ID int) string {
	return fmt.Sprintf("person-%d", ID)
}

func (d *DAO) logger() logging.Logger {
	return d.cfg.Logger()
}
