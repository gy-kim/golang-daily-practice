package mux

import (
	"net/http"
)

type Route struct {
	handler http.Handler

	buildOnly bool

	name string

	err error

	namedRoutes map[string]*Route

	routeConf
}
