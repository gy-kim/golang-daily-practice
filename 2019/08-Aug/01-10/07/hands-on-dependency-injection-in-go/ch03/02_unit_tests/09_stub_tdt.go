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
		dest         string
		loaderStub   *PersonLoaderStub
		expectedName string
		expectErr    bool
	}{
		{
			dest: "happ path",
			loaderStub: &PersonLoaderStub{
				Person: &Person{Name: "Sophia"},
			},
			expectedName: "Sophoa",
			expectErr:    false,
		},
		{
			dest: "input error",
			loaderStub: &PersonLoaderStub{
				Error: ErrNotFound,
			},
			expectedName: "",
			expectErr:    true,
		},
		{
			dest: "system error path",
			loaderStub: &PersonLoaderStub{
				Error: errors.New("somthing failed"),
			},
			expectedName: "",
			expectErr:    true,
		},
	}

	for _, s := range scenarios {
		result, resultErr := LoadPersonName(s.loaderStub, fakeID)

		assert.Equal(t, s.expectedName, result, s.dest)
		assert.Equal(t, s.expectErr, resultErr != nil, s.dest)
	}
}
