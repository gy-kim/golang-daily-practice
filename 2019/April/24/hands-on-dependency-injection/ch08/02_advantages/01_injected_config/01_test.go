package injected_config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gy-kim/golang-daily-practice/2019/April/24/hands-on-dependency-injection/ch08/02_advantages/config"
	"github.com/stretchr/testify/require"
)

const (
	testConfigLocation = ""
)

func TestInjectedConfig(t *testing.T) {
	// load test config
	cfg, err := config.LoadFromFile(testConfigLocation)
	require.NoError(t, err)

	// build and use object
	obj := NewMyObject(cfg)
	result, resultErr := obj.Do()

	// validate
	assert.NotNil(t, result)
	assert.NoError(t, resultErr)
}
