package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "unauthorized - missing token", http.StatusUnauthorized)
			return
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "super_secret_key_for_development"
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "unauthorized - invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "unauthorized - invalid claims", http.StatusUnauthorized)
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "unauthorized - invalid subject", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(sub)
		if err != nil {
			http.Error(w, "unauthorized - invalid user id format", http.StatusUnauthorized)
			return
		}

		user := &User{
			ID: userID,
		}

		if roleStr, ok := claims["role"].(string); ok {
			user.Role = roleStr
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value("user").(*User)
			if !ok {
				http.Error(w, "unauthorized - no user context", http.StatusUnauthorized)
				return
			}

			if user.Role != requiredRole {
				http.Error(w, "forbidden - insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
