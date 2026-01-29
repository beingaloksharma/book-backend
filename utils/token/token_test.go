package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestTokenFlow(t *testing.T) {
	// Setup
	viper.Set("jwt.secret", "testsecret")
	Init()

	userID := uint(123)
	role := "admin"

	// 1. Generate Token
	tokenString, err := GenerateToken(userID, role)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// 2. Validate Token
	claims, err := ValidateToken(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// Check Claims
	// jwt.MapClaims unmarshals numbers as float64
	uid := uint(claims["user_id"].(float64))
	r := claims["role"].(string)

	assert.Equal(t, userID, uid)
	assert.Equal(t, role, r)
}

func TestInvalidToken(t *testing.T) {
	viper.Set("jwt.secret", "testsecret")
	Init()

	_, err := ValidateToken("invalid-token-string")
	assert.Error(t, err)
}

func TestExpiredToken(t *testing.T) {
	viper.Set("jwt.secret", "testsecret")
	Init()

	// Manually create an expired token
	claims := jwt.MapClaims{
		"user_id": 1,
		"role":    "user",
		"exp":     time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("testsecret"))

	// Validate
	_, err := ValidateToken(tokenString)
	assert.Error(t, err)
}
