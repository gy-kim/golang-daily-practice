package exchange

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/21-31/26-31/hands-on-dependency-injection-in-go/ch09/acme/internal/logging"
	"github.com/stretchr/testify/assert"
)

func TestInternalBoundaryTest(t *testing.T) {
	// start our test server
	server := httptest.NewServer(&happyExchangeRateService{})
	defer server.Close()

	// define the config
	cfg := &testConfig{
		baseURL: server.URL,
		apiKey:  "",
	}

	// crate a converter to test
	converter := NewConverter(cfg)
	resultRate, resultErr := converter.Exchange(context.Background(), 100.00, "AUD")

	// validate the result
	assert.Equal(t, 101.01, resultRate)
	assert.NoError(t, resultErr)
}

type happyExchangeRateService struct{}

func (*happyExchangeRateService) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	payload := []byte(`
{
   "success":true,
   "historical":true,
   "date":"2010-11-09",
   "timestamp":1289347199,
   "source":"USD",
   "quotes":{
      "USDAUD":0.989981
   }
}`)
	response.Write(payload)
}

type testConfig struct {
	baseURL string
	apiKey  string
}

// Logger implements Config
func (t *testConfig) Logger() logging.Logger {
	return &logging.LoggerStdOut{}
}

// ExchangeBaseURL implements Config
func (t *testConfig) ExchangeBaseURL() string {
	return t.baseURL
}

func (t *testConfig) ExchangeAPIKey() string {
	return t.apiKey
}
