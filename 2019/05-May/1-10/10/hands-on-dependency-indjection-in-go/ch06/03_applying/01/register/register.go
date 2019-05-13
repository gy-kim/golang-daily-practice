package register

import "github.com/gy-kim/golang-daily-practice/2019/05-May/1-10/10/hands-on-dependency-indjection-in-go/ch06/03_applying/01/data"

type Registerer struct {
}

func (r *Registerer) Do(in *data.Person) (int, error) {
	return 0, nil
}
