package jwt

import "errors"

// Errors.
var (
	// ErrNotFoundFromContext an error when couldn't find a token in context.
	ErrNotFoundFromContext = errors.New("couldn't find a token in context")

	// ErrNotFoundFromCookie an error when couldn't find a token in cookie.
	ErrNotFoundFromCookie = errors.New("couldn't find a token in cookie")
)
