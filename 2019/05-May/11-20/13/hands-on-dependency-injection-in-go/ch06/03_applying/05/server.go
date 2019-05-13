package reset

import (
	"github.com/gy-kim/golang-daily-practice/2019/05-May/11-20/13/hands-on-dependency-injection-in-go/ch06/03_applying/05/get"
	"github.com/gy-kim/golang-daily-practice/2019/05-May/11-20/13/hands-on-dependency-injection-in-go/ch06/03_applying/05/list"
	"github.com/gy-kim/golang-daily-practice/2019/05-May/11-20/13/hands-on-dependency-injection-in-go/ch06/03_applying/05/register"
)

func New(address string) *Server {
	return &Server{
		address:         address,
		handlerGet:      NewGetHandler(&get.Getter{}),
		handlerList:     NewListHandler(&list.Lister{}),
		handlerNotFound: notFoundHandler,
		handlerRegister: NewRegisterHandler(&register.Registerer{}),
	}
}
