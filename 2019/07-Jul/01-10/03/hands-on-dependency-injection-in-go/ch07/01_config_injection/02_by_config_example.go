package config_injection

import "time"

func NewByConfigConstructor(cfg MyConfig, limiter RateLimiter, cache Cache) *MyStruct {
	return &MyStruct{}
}

// MyConfig defines the config the MyStruct
type MyConfig interface {
	Logger() Logger
	Instrumentation() Instrumentation
	Timeout() time.Duration
	Workers() int
}
