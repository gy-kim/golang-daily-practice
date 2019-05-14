package injected_config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gy-kim/golang-daily-practice/2019/05-May/11-20/14/hands-on-dependency-injection-in-go/ch08/02_advantages/config"
)

const (
	testConfigLocation = ""
)

func TestInjectedConfig(t *testing.T) {
	// laod test config
	cfg, err := config.LoadFromFile(testConfigLocation)
	require.NoError(t, err)

	// build and use object
	obj := NewMyObject(cfg)
	result, resultErr := obj.Do()

	// vaidate
	assert.NotNil(t, result)
	assert.NoError(t, resultErr)
}
