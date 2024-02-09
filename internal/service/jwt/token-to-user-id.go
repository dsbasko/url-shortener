package jwt

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/entities"
)

// TokenToUserID parses jwt token and returns user ID.
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
