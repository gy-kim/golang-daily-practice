package config_injection

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gy-kim/golang-daily-practice/2019/April/24/hands-on-dependency-injection/ch08/02_advantages/logging"
	"github.com/gy-kim/golang-daily-practice/2019/April/24/hands-on-dependency-injection/ch08/02_advantages/stats"
)

func TestConfigInjection(t *testing.T) {
	// build test config
	cfg := &TestConfig{}

	// build and use object
	obj := NewMyObject(cfg)
	result, resultErr := obj.Do()

	// validate
	assert.NotNil(t, result)
	assert.NoError(t, resultErr)
}

type TestConfig struct {
	logger *logging.Logger
	stats  *stats.Collector
}

func (t *TestConfig) Logger() *logging.Logger {
	return t.logger
}

func (t *TestConfig) Stats() *stats.Collector {
	return t.stats
}
