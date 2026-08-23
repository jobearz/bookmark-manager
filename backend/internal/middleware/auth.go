package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jobearz/bookmark-manager/config"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get token from auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// parse and identify token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			secret := config.JWTSecret()
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		// if valid, call next (w, r) to the actual handler
		next(w, r)
	}
}

func GetUserIDFromToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		secret := config.JWTSecret()
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return ""
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims["user_id"].(string)
}
