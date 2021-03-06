package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gy-kim/golang-daily-practice/2019/05-May/1-10/02/hands-on-dependency-injection/ch09/acme/internal/logging"
)

// Config is the config for the REST package
type Config interface {
	Logger() logging.Logger
	BindAddress() string
}

// New will create and initialize the server
func New(cfg Config,
	getModel GetModel,
	listModel ListModel,
	registerModel RegisterModel) *Server {

	return &Server{
		address:         cfg.BindAddress(),
		handlerGet:      NewGetHandler(cfg, getModel),
		handlerList:     NewListHandler(listModel),
		handlerNotFound: notFoundHandler,
		handlerRegister: NewRegisterHandler(registerModel),
	}
}

// Server is the HTTP REST server
type Server struct {
	address string
	server  *http.Server

	handlerGet      http.Handler
	handlerList     http.Handler
	handlerNotFound http.HandlerFunc
	handlerRegister http.Handler
}

// Listen will start a HTTP rest for this service
func (s *Server) Listen(stop <-chan struct{}) {
	router := s.buildRouter()

	// create the HTTP server
	s.server = &http.Server{
		Handler: router,
		Addr:    s.address,
	}

	// listen for shutdown
	go func() {
		// wait for shutdown signal
		<-stop

		_ = s.server.Close()
	}()

	_ = s.server.ListenAndServe()
}

func (s *Server) buildRouter() http.Handler {
	router := mux.NewRouter()

	// map URL endpoints to HTTP handlers
	router.Handle("/person/{id}/", s.handlerGet).Methods("GET")
	router.Handle("/person/list", s.handlerList).Methods("GET")
	router.Handle("/person/register", s.handlerRegister).Methods("POST")

	// convert a "catch all" not found handler
	router.NotFoundHandler = s.handlerNotFound

	return router
}
