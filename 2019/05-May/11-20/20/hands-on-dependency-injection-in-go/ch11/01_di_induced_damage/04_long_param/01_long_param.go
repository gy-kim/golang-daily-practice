package long_param

import (
	"net/http"
	"time"
)

// MyHandler does something fantastic
type MyHandler struct {
	config    Config
	parser    Parser
	formatter Formatter
	limiter   RateLimiter
	loader    Loader
}

func (m *MyHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ID, err := m.parser.Extract(request)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := m.loader.Load(ID)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = m.formatter.Format(response, data)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
	}
}

// Config combines environmental concerns like logging and instrumentation with any other config
type Config interface {
	Logger() Logger
	Instrumentation() Instrumentation
}

// Logger logs stuff
type Logger interface {
	Error(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Info(message string, args ...interface{})
	Debug(message string, args ...interface{})
}

// Instrumentation records the performances and events
type Instrumentation interface {
	Count(key string, value int)
	Duration(key string, start time.Time)
}

// Parser will extract detailed from the request
type Parser interface {
	Extract(req *http.Request) (int, error)
}

// Formatter will build the output
type Formatter interface {
	Format(resp http.ResponseWriter, data []byte) error
}

// FancyFormatter Implements Formatter
type FancyFormatter struct{}

func (f *FancyFormatter) Format(response http.ResponseWriter, data []byte) error {
	// does something fancy with the data
	_, err := response.Write([]byte(`something fancy!`))
	return err
}

// RateLimiter limits how many concurrency requests we can make or process
type RateLimiter interface {
	Acquire()
	Release()
}

type Loader interface {
	Load(ID int) ([]byte, error)
}
