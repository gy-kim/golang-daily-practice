package list

import "github.com/gy-kim/golang-daily-practice/2019/05-May/11-20/13/hands-on-dependency-injection-in-go/ch06/03_applying/05/data"

type Lister struct{}

func (l *Lister) Do() ([]*data.Person, error) {
	return nil, nil
}
