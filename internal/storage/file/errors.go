package file

import "fmt"

var (
	ErrURLNotFound       = fmt.Errorf("url not found")
	ErrURLsNotFound      = fmt.Errorf("urls not found")
	ErrURLFoundDuplicate = fmt.Errorf("found duplicate url")
)
