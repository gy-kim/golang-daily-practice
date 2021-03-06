package unit_tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func concat(a, b string) string {
	return a + b
}

func TestTooSimple(t *testing.T) {
	a := "hello "
	b := "world"
	expected := "hello world"

	assert.Equal(t, expected, concat(a, b))
}
