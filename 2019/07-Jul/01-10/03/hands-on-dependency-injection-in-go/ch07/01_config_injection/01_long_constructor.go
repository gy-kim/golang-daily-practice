package config_injection

import "time"

func NewLongConstructor(logger Logger, stats Instrumentation, limiter RateLimiter, cache Cache, timeout time.Duration, workers int) *MyStruct {
	return &MyStruct{}
}

// MyStruct does something fantastic
type MyStruct struct{}

// Logger logs stuff
type Logger interface {
	Error(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Info(message string, args ...interface{})
	Debug(message string, args ...interface{})
}

type Instrumentation interface {
	Count(key string, value int)
	Duration(key string, start time.Time)
}

type RateLimiter interface {
	Acquire()
	Release()
}

type Cache interface {
	Store(key string, data []byte)
	Get(key string) ([]byte, error)
}
