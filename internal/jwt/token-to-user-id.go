package jwt

import (
	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/entities"
	"github.com/golang-jwt/jwt/v5"
)

func TokenToUserID(tokenString string) string {
	claims := &entities.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return config.GetJWTSecret(), nil
	})
	if err != nil || !token.Valid {
		return ""
	}

	return claims.UserID
}
