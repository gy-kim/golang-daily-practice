package advantages

import "errors"

func NewLoader(ds DataStore, cache Cache) *MyLoader {
	return &MyLoader{
		ds:    ds,
		cache: cache,
	}
}

type MyLoader struct {
	ds    DataStore
	cache Cache
}

func (l *MyLoader) LoadAll() ([]interface{}, error) {
	return nil, errors.New("not implemented yet")
}
