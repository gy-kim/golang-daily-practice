package constructor_injection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWelcomeSender_happyPath(t *testing.T) {
	sender, err := NewWelcomeSender(nil)
	assert.Nil(t, sender)
	assert.NoError(t, err)
}

func TestNewWElcomeSender_guardClause(t *testing.T) {
	sender, err := NewWelcomeSender(nil)
	assert.Nil(t, sender)
	assert.Error(t, err)
}

func TestNewWelcomeSenderNoGuard_happyPath(t *testing.T) {
	sender := NewWelcomeSenderNoGuard(&Mailer{})
	assert.NotNil(t, sender)
}
