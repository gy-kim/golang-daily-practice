package get

import (
	"errors"

	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/01-10/08-09/hands-on-dependency-injection-in-go/ch04/acme/internal/modules/data"
)

var (
	errPersonNotFound = errors.New("person not found")
)

type Getter struct{}

func (g *Getter) Do(ID int) (*data.Person, error) {
	// load person from the data layer
	person, err := data.Load(ID)
	if err != nil {
		if err == data.ErrNotFound {
			return nil, errPersonNotFound
		}
		return nil, err
	}
	return person, err
}
