package testutils

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateTestToken creates a valid JWT token cookie for testing.
func GenerateTestToken(userID string, role string) *http.Cookie {
	secret := "super_secret_key_for_development" // default from middleware

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(time.Hour).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(secret))

	return &http.Cookie{
		Name:  "token",
		Value: tokenString,
	}
}
