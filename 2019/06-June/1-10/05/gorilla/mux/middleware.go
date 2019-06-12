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
