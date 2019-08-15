package main

import (
	"fmt"
	"io"
	"net/http"
)

type MyServer struct{}

func (m *MyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello Decorator!")
}

type LoggerMiddleware struct {
	Handler   http.Handler
	LogWriter io.Writer
}

func (l *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(l.LogWriter, "Request URL: %s\n", r.RequestURI)
}

func main() {}
