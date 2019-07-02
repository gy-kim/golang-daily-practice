package list

import "github.com/gy-kim/golang-daily-practice/2019/07-Jul/01-10/02/hands-on-dependency-injection-in-go/03_applying/data"

type Lister struct{}

func (l *Lister) Do() ([]*data.Person, error) {
	return nil, nil
}
