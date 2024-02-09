package entity

import "github.com/golang-jwt/jwt/v5"

// JWTClaims is a jwt claims.
type JWTClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}
