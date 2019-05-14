package config_injection

import "time"

func NewConfigConstructor(cfg MyConfig, limiter RateLimiter, cache *Cache) *MyStruct {
	return &MyStruct{}
}

type MyConfig interface {
	Logger() Logger
	Instrumentation() Instrumentation
	Timeout() time.Duration
	Workers() int
}
