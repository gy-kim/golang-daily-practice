package get

import "github.com/gy-kim/golang-daily-practice/2019/05-May/11-20/13/hands-on-dependency-injection-in-go/ch06/03_applying/05/data"

type Getter struct{}

func (g *Getter) Do(ID int) (*data.Person, error) {
	return nil, nil
}
