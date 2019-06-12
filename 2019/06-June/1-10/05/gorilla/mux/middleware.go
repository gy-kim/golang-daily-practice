package mux

import "net/http"

type MiddlewareFunc func(http.Handler) http.Handler

type middleware interface {
	Middleware(handler http.Handler) http.Handler
}

func (mw MiddlewareFunc) Middleware(handler http.Handler) http.Handler {
	return mw(handler)
}

func (r *Router) Use(mwf ...MiddlewareFunc) {
	for _, fn := range mwf {
		r.middlewares = append(r.middlewares, fn)
	}
}

func (r *Router) useInterface(mw middleware) {
	r.middlewares = append(r.middlewares, mw)
}

func CORSMethodMiddleware(r *Router) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		var allMethods []string

		err := r.Walk(func(route *Route, _ *Router, _ []*Route) error {

		})
	}
}
