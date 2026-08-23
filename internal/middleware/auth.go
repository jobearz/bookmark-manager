package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jobearz/bookmark-manager/config"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		tokenString := cookie.Value
		secret := config.JWTSecret()
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetUserIDFromToken(r *http.Request) string {
	cookie, err := r.Cookie("token")
	if err != nil {
		return ""
	}
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret()), nil
	})
	if err != nil || !token.Valid {
		return ""
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims["user_id"].(string)
}
