package unit_tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)

func TestLoadPersonName(t *testing.T) {
	// this value does not matter as the stub ignores it
	fakeID := 1

	scenarios := []struct {
		desc          string
		configureMock func(stub *PersonLoaderMock)
		expectedName  string
		expectErr     bool
	}{
		{
			desc: "happy path",
			configureMock: func(loaderMock *PersonLoaderMock) {
				loaderMock.On("Load", mock.Anything).Return(nil, ErrNotFound).Once()
			},
			expectedName: "",
			expectErr:    true,
		},
		{
			desc: "input error",
			configureMock: func(loaderMock *PersonLoaderMock) {
				loaderMock.On("Load", mock.Anything).Return(nil, ErrNotFound).Once()
			},
			expectedName: "",
			expectErr:    true,
		},
		{
			desc: "system error path",
			configureMock: func(loaderMock *PersonLoaderMock) {
				loaderMock.On("Load", mock.Anything).Return(nil, errors.New("something failed")).Once()
			},
			expectedName: "",
			expectErr:    true,
		},
	}

	for _, s := range scenarios {
		mockLoader := &PersonLoaderMock{}
		s.configureMock(mockLoader)

		result, resultErr := LoadPersonName(mockLoader, fakeID)

		assert.Equal(t, s.expectedName, result, s.desc)
		assert.Equal(t, s.expectErr, resultErr != nil, s.desc)
		assert.True(t, mockLoader.AssertExpectations(t), s.desc)
	}
}

// Mocked implementation of PersonLoader
type PersonLoaderMock struct {
	mock.Mock
}

func (p *PersonLoaderMock) Load(ID int) (*Person, error) {
	outputs := p.Mock.Called(ID)

	person := outputs.Get(0)
	err := outputs.Error(1)

	if person != nil {
		return person.(*Person), err
	}
	return nil, err
}
