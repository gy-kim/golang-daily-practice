package register

import (
	context "context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/mock"

	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/21-31/26-31/hands-on-dependency-injection-in-go/ch09/acme/internal/logging"
	data "github.com/gy-kim/golang-daily-practice/2019/08-Aug/21-31/26-31/hands-on-dependency-injection-in-go/ch09/acme/internal/modules/data"
)

func TestRegisterer_Do_happyPath(t *testing.T) {
	// configure the mock saver
	mockResult := 888

	mockSaver := &mockMySaver{}
	mockSaver.On("Save", mock.Anything, mock.Anything).Return(mockResult, nil).Once()

	// define context and therefore test timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// inputs
	in := &data.Person{
		FullName: "Chang",
		Phone:    "11122233355",
		Currency: "CNY",
	}

	// call method
	registerer := &Registerer{
		cfg:       &testConfig{},
		exchanger: &stubExchange{},
		data:      mockSaver,
	}
	ID, err := registerer.Do(ctx, in)

	// validate expectations
	require.NoError(t, err)
	assert.Equal(t, 888, ID)
	assert.True(t, mockSaver.AssertExpectations(t))
}

func TestRegisterer_Do_error(t *testing.T) {
	// configure the mock saver
	mockSaver := &mockMySaver{}
	mockSaver.On("Save", mock.Anything, mock.Anything).Return(0, errors.New("something failed"))

	// define context and therefore test timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// inputs
	in := &data.Person{
		FullName: "Chang",
		Phone:    "11122233355",
		Currency: "CNY",
	}

	// call method
	registerer := &Registerer{
		cfg:       &testConfig{},
		exchanger: &stubExchange{},
		data:      mockSaver,
	}
	ID, err := registerer.Do(ctx, in)

	// validate expectations
	require.Error(t, err)
	assert.Equal(t, 0, ID)
	assert.True(t, mockSaver.AssertExpectations(t))
}

// Stub implementation of Config
type testConfig struct{}

// Logger implement Config
func (t *testConfig) Logger() logging.Logger {
	return &logging.LoggerStdOut{}
}

// RegistrationBasePrice implement Config
func (t *testConfig) RegistrationBasePrice() float64 {
	return 12.34
}

// DataDSN implements Config
func (t *testConfig) DataDSN() string {
	return ""
}

type stubExchange struct{}

func (s stubExchange) Exchange(ctx context.Context, basePrice float64, currency string) (float64, error) {
	return 12.34, nil
}
