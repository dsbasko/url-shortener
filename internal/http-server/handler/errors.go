package handler

import "errors"

var (
	ErrEmptyBody        = errors.New("empty body")
	ErrWrongContentType = errors.New("wrong content type")
)
