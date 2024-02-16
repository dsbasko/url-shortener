package jwt

import (
	"fmt"
	"net/http"
)

// GetFromCookie gets jwt token from cookie.
func GetFromCookie(r *http.Request) (string, error) {
	token, err := r.Cookie(CookieKey)

	if err != nil {
		return "", fmt.Errorf("failed to get cookie value: %w", err)
	}

	if token == nil {
		return "", ErrNotFoundFromCookie
	}

	return token.Value, nil
}