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

func (r *Route) SkipClean() bool {
	return r.skipClean
}

func (r *Route) Match(req *http.Request, match *RouteMatch) bool {
	if r.buildOnly || r.err != nil {
		return false
	}

	var matchErr error

	for _, m := range r.matchers {
		if matched := m.Match(req, match); !matched {
			if _, ok := m.(methodMatcher); ok {
				matchErr = ErrMethodMismatch
				continue
			}

			if match.MatchErr == ErrNotFound {
				match.MatchErr = nil
			}

			matchErr = nil
			return false
		}
	}

	if matchErr != nil {
		match.MatchErr = matchErr
		return false
	}

	if match.MatchErr == ErrMethodMismatch {
		match.MatchErr = nil
		match.Handler = r.handler
	}

	if match.Route == nil {
		match.Route = r
	}
	if match.Handler == nil {
		match.Handler = r.handler
	}
	if match.Vars == nil {
		match.Vars = make(map[string]string)
	}

	r.regexp.setMatch(req, match, r)
	return true
}

type matcher interface {
	Match(*http.Request, *RouteMatch) bool
}

type BuildVarsFunc func(map[string]string) map[string]string

type methodMatcher []string

func (m methodMatcher) Match(r *http.Request, match *RouteMatch) bool {
	return matchInArray(m, r.Method)
}
