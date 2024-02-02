package errors

import (
	goErrors "errors"
	"fmt"
)

func Example() {
	// Crete a new error
	errFirst := goErrors.New("first error")

	// Wrap the error
	errSecond := fmt.Errorf("second error: %w", errFirst)
	errThird := fmt.Errorf("third error: %w", errSecond)

	err := UnwrapAll(errThird)
	fmt.Println(err)
	// Output: first error
}
