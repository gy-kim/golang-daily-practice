package get

import (
	"errors"
	"testing"

	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/21-31/26-31/hands-on-dependency-injection-in-go/ch09/acme/internal/modules/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetter_Do_happyPath(t *testing.T) {
	// inputs
	ID := 1234

	// configure the mock loader
	mockResult := &data.Person{
		ID:       1234,
		FullName: "Doug",
	}
	mockLoader := &mockMyLoader{}
	mockLoader.On("Load", mock.Anything, ID).Return(mockResult, nil).Once()

	// call method
	getter := &Getter{
		data: mockLoader,
	}
	person, err := getter.Do(ID)

	// validate expecteations
	require.NoError(t, err)
	assert.Equal(t, ID, person.ID)
	assert.Equal(t, "Doug", person.FullName)
	assert.True(t, mockLoader.AssertExpectations(t))
}

func TestGetter_Do_noSuchPerson(t *testing.T) {
	// inputs
	ID := 5678

	// configure the mock loader
	mockLoader := &mockMyLoader{}
	mockLoader.On("Load", mock.Anything, ID).Return(nil, data.ErrNotFound).Once()

	// call method
	getter := &Getter{
		data: mockLoader,
	}
	person, err := getter.Do(ID)

	// validate expectations
	require.Equal(t, errPersonNotFound, err)
	assert.Nil(t, person)
	assert.True(t, mockLoader.AssertExpectations(t))
}

func TestGetter_Do_error(t *testing.T) {
	// inputs
	ID := 1234

	// configure the mock loader
	mockLoader := &mockMyLoader{}
	mockLoader.On("Load", mock.Anything, ID).Return(nil, errors.New("something failed"))

	// call method
	getter := &Getter{
		data: mockLoader,
	}
	person, err := getter.Do(ID)

	// validate expectations
	require.Error(t, err)
	assert.Nil(t, person)
	assert.True(t, mockLoader.AssertExpectations(t))
}
