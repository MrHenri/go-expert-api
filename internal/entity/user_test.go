package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	u, err := NewUser("John", "John@email.com", "12345")
	assert.Nil(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, "John", u.Name)
	assert.Equal(t, "John@email.com", u.Email)
	assert.NotEmpty(t, u.ID)
	assert.NotEmpty(t, u.Password)
}

func TestValidatePassword(t *testing.T) {
	u, err := NewUser("John", "John@email.com", "12345")
	assert.Nil(t, err)
	assert.NotNil(t, u)
	assert.True(t, u.ValidatePassword("12345"))
	assert.False(t, u.ValidatePassword("123456"))
	assert.NotEqual(t, u.Password, "12345")
}
