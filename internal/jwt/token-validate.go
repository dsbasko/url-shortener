package jwt

import (
	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func TokenValidate(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return config.GetJWTSecret(), nil
	})

	if err != nil || token == nil {
		return false
	}

	return token.Valid
}
