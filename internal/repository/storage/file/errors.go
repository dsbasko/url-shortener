package file

import "fmt"

// Errors.
var (
	// ErrURLNotFound an error returned when a URL is not found.
	ErrURLNotFound = fmt.Errorf("url not found")

	// ErrURLsNotFound an error returned when URLs are not found.
	ErrURLsNotFound = fmt.Errorf("urls not found")

	// ErrURLFoundDuplicate an error returned when a duplicate URL is found.
	ErrURLFoundDuplicate = fmt.Errorf("found duplicate url")
)
