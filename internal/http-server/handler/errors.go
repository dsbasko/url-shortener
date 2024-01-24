package handler

import "errors"

// Errors.
var (
	// ErrEmptyBody returned when the body is empty.
	ErrEmptyBody = errors.New("empty body")

	// ErrWrongContentType returned when the content type is wrong.
	ErrWrongContentType = errors.New("wrong content type")
)
