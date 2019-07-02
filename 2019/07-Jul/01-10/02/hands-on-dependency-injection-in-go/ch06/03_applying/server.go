package reset

import (
	"github.com/gy-kim/golang-daily-practice/2019/07-Jul/01-10/02/hands-on-dependency-injection-in-go/03_applying/get"
	"github.com/gy-kim/golang-daily-practice/2019/07-Jul/01-10/02/hands-on-dependency-injection-in-go/03_applying/list"
	"github.com/gy-kim/golang-daily-practice/2019/07-Jul/01-10/02/hands-on-dependency-injection-in-go/03_applying/register"
)

func New(address string) *Server {
	return &Server{
		address:         address,
		handlerGet:      NewGetHandler(&get.Getter{}),
		handlerList:     NewListHandler(&list.Lister{}),
		handlerNotFound: noFoundHandler,
		handlerRegister: NewRegisterHandler(&register.Registerer{}),
	}
}
