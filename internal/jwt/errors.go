package jwt

import "errors"

var (
	ErrNotFoundFromContext = errors.New("couldn't find a token in context")
	ErrNotFoundFromCookie  = errors.New("couldn't find a token in cookie")
)
