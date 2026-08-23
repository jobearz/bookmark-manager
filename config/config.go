package config

import "os"

func JWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "dev-secret-change-in-production"
	}
	return secret
}
