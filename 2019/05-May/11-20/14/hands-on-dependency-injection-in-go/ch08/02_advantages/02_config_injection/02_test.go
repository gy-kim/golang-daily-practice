package config_injection

import (
	"testing"

	"github.com/gy-kim/golang-daily-practice/2019/05-May/11-20/14/hands-on-dependency-injection-in-go/ch08/02_advantages/logging"
	"github.com/gy-kim/golang-daily-practice/2019/05-May/11-20/14/hands-on-dependency-injection-in-go/ch08/02_advantages/stats"
	"github.com/stretchr/testify/assert"
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
