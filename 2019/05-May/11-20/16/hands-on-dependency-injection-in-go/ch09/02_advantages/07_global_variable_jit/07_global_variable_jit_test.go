package global_variable_jit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSavor_Do(t *testing.T) {
	// input
	carol := &User{
		Name:     "Carol",
		Password: "IamKing",
	}

	// mocks/stubs
	stubStorage := &StubUserStorage{}

	// do call
	savor := &Savor{
		storage: stubStorage,
	}
	resultErr := savor.Do(carol)

	// valiate
	assert.NotEqual(t, resultErr, "unexpected error")
}

type StubUserStorage struct{}

func (s *StubUserStorage) Save(_ *User) error {
	// return "happy path"
	return nil
}
