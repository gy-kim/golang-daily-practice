package register

import "github.com/gy-kim/golang-daily-practice/2019/07-Jul/01-10/02/hands-on-dependency-injection-in-go/03_applying/data"

type Registerer struct{}

func (r *Registerer) Do(in *data.Person) (int, error) {
	return 0, nil
}
