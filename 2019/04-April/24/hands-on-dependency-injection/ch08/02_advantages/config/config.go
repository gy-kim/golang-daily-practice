package config

import (
	"sync"

	"github.com/gy-kim/golang-daily-practice/2019/April/24/hands-on-dependency-injection/ch08/02_advantages/stats"

	"github.com/gy-kim/golang-daily-practice/2019/April/24/hands-on-dependency-injection/ch08/02_advantages/logging"
)

func LoadFromFile(path string) (*Config, error) {
	// TODO: implement
	return &Config{}, nil
}

type Config struct {
	// Log config
	LogLevel       int `json:"log_level"`
	logger         *logging.Logger
	loggerInitOnce sync.Once

	// Instrumentation config
	StatsDHostAndPort string `json:"stats_d_host_and_port"`
	stats             *stats.Collector
	statsInitOnce     sync.Once

	// Rate Limiter config
	RateLimiterMaxConcurrent int `json:"rate_limiter_max_concurrent"`
}

func (c *Config) Logger() *logging.Logger {
	c.loggerInitOnce.Do(func() {
		// use log level to create new logger
		c.logger = &logging.Logger{
			Level: c.LogLevel,
		}
	})
	return c.logger
}

func (c *Config) Stats() *stats.Collector {
	c.statsInitOnce.Do(func() {
		c.stats = &stats.Collector{
			HostAndPort: c.StatsDHostAndPort,
		}
	})
	return c.stats
}
