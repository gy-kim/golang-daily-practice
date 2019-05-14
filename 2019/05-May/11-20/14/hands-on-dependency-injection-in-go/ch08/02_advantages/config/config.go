package config

import (
	"sync"

	"github.com/gy-kim/golang-daily-practice/2019/05-May/11-20/14/hands-on-dependency-injection-in-go/ch08/02_advantages/logging"
	"github.com/gy-kim/golang-daily-practice/2019/05-May/11-20/14/hands-on-dependency-injection-in-go/ch08/02_advantages/stats"
)

// LoadFromFile loads the config from the supplied path
func LoadFromFile(path string) (*Config, error) {
	return &Config{}, nil
}

// Config is the result of loading config from a file
type Config struct {
	// Log config
	LogLevel       int `json:"log_level"`
	logger         *logging.Logger
	loggerInitOnce sync.Once

	// Instrumentation config
	StatesDHostAndPort string `json:"stats_d_host_and_port"`
	stats              *stats.Collector
	statsInitOnce      sync.Once
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
			HostAndPort: c.StatesDHostAndPort,
		}
	})

	return c.stats
}
