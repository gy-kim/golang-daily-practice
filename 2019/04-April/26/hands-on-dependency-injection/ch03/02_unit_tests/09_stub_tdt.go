package unit_tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadPersonNameStubs(t *testing.T) {
	// this value does not matter as the stub ignores it
	fakeID := 1

	scenarios := []struct {
		desc         string
		loaderStub   *PersonLoaderStub
		expectedName string
		expectedErr  bool
	}{
		{
			desc: "happy path",
			loaderStub: &PersonLoaderStub{
				Person: &Person{Name: "Sophia"},
			},
			expectedName: "Sophia",
			expectedErr:  false,
		}, {
			desc: "input error",
			loaderStub: &PersonLoaderStub{
				Error: ErrNotFound,
			},
			expectedName: "",
			expectedErr:  true,
		}, {
			desc: "system error path",
			loaderStub: &PersonLoaderStub{
				Error: errors.New("something failed"),
			},
			expectedName: "",
			expectedErr:  false,
		},
	}

	for _, scenario := range scenarios {
		result, resultErr := LoadPersonName(scenario.loaderStub, fakeID)

		assert.Equal(t, scenario.expectedName, result, scenario.desc)
		assert.Equal(t, scenario.expectedErr, resultErr != nil, scenario.desc)

	}
}
