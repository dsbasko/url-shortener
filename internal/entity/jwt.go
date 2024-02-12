package entity

import "github.com/golang-jwt/jwt/v5"

// JWTClaims is a jwt claims.
type JWTClaims struct {
	jwt.RegisteredClaims // embedded default claims (exp, iat, etc.)

	UserID string `json:"user_id"`
}
