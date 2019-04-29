package list

import (
	"errors"

	"github.com/gy-kim/golang-daily-practice/2019/April/29/hands-on-dependency-injection/ch04/acme/internal/modules/data"
)

var (
	// error thrown when there are no people in the database
	errPeopleNotFound = errors.New("no people found")
)

// Lister will attempt to load all people in the database
// It can return an error caused by the data layer.
type Lister struct{}

func (l *Lister) load() ([]*data.Person, error) {
	people, err := data.LoadAll()
	if err != nil {
		if err == data.ErrNotFound {
			// by converting the error we are encapsulating the implementation details from our users.
			return nil, errPeopleNotFound
		}
		return nil, err
	}
	return people, nil
}
