// +build external

package exchange

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gy-kim/golang-daily-practice/2019/05-May/1-10/02/hands-on-dependency-injection/ch09/acme/internal/config"
)

func TestExternalBoundaryTest(t *testing.T) {
	// define the config
	cfg := &testConfig{
		baseURL: config.App.ExchangeRateBaseURL,
		apiKey:  config.App.ExchangeRateAPIKey,
	}

	// create a converter to test
	converter := NewConverter(cfg)

	// fetch from the server
	response, err := converter.loadRateFromServer(context.Background(), "AUD")
	require.NotNil(t, response)
	require.NoError(t, err)

	// parse the response
	resultRate, err := converter.extractRate(response, "AUD")
	require.NoError(t, err)

	// validate the result
	assert.True(t, resultRate > 0)
}
