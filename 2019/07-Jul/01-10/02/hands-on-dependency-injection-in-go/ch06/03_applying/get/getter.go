package get

import "github.com/gy-kim/golang-daily-practice/2019/07-Jul/01-10/02/hands-on-dependency-injection-in-go/03_applying/data"

// Stub implementation so that the example comiles
type Getter struct{}

func (g *Getter) Do(ID int) (*data.Person, error) {
	return nil, nil
}
