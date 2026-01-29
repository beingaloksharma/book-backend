package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var secretKey []byte

func Init() {
	secret := viper.GetString("jwt.secret")
	if secret == "" {
		secret = "supersecretkey" // Fallback should be changed in production
	}
	secretKey = []byte(secret)
}

func GenerateToken(userID uint, role string) (string, error) {
	if len(secretKey) == 0 {
		Init()
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	if len(secretKey) == 0 {
		Init()
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
