package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "mySecretPassword"
	hash, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mySecretPassword"
	hash, _ := HashPassword(password)

	match := CheckPasswordHash(password, hash)
	assert.True(t, match, "Password should match hash")

	match = CheckPasswordHash("wrongPassword", hash)
	assert.False(t, match, "Wrong password should not match hash")
}
