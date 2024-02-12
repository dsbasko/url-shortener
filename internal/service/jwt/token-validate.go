package jwt

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
)

// TokenValidate validates jwt token.
func TokenValidate(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return config.JWTSecret(), nil
	})

	if err != nil || token == nil {
		return false
	}

	return token.Valid
}
