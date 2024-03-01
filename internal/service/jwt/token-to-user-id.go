package jwt

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/dsbasko/url-shortener/internal/config"
	"github.com/dsbasko/url-shortener/internal/entity"
)

// TokenToUserID parses jwt token and returns user ID.
func TokenToUserID(tokenString string) string {
	claims := &entity.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return config.JWTSecret(), nil
	})
	if err != nil || !token.Valid {
		return ""
	}

	return claims.UserID
}
