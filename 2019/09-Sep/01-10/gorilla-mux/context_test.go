package mux

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestNativeContextMiddleware(t *testing.T) {
	withTimeout := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
			defer cancel()
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	r := NewRouter()
	r.Handle("/path/{foo}", withTimeout(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := Vars(r)
		if vars["foo"] != "bar" {
			t.Fatal("Expected foo var to be set")
		}
	})))

}
