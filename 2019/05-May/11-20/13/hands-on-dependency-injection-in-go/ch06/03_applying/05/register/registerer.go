package register

import "github.com/gy-kim/golang-daily-practice/2019/05-May/11-20/13/hands-on-dependency-injection-in-go/ch06/03_applying/05/data"

type Registerer struct{}

func (r *Registerer) Do(in *data.Person) (int, error) {
	return 0, nil
}
