package get

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetter_Do(t *testing.T) {
	// inputs
	ID := 1
	name := "John"

	// calll method
	getter := &Getter{}
	person, err := getter.Do(ID)

	require.NoError(t, err)
	assert.Equal(t, ID, person.ID)
	assert.Equal(t, name, person.FullName)
}
