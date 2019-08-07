package unit_tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRound(t *testing.T) {
	scenarios := []struct {
		dest     string
		in       float64
		expected int
	}{
		{
			dest:     "round down",
			in:       1.1,
			expected: 1,
		},
		{
			dest:     "round up",
			in:       3.7,
			expected: 4,
		},
		{
			dest:     "unchanged",
			in:       6.0,
			expected: 6,
		},
	}

	for _, s := range scenarios {
		in := float64(s.in)

		result := Round(in)
		assert.Equal(t, s.expected, result)
	}
}
