package urls

import (
	"math/rand"
	"time"
)

func RandomString(length int) string {
	if length < 1 {
		return ""
	}
	rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	characters := `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789`
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = characters[rand.Intn(len(characters))] //nolint:gosec
	}
	return string(result)
}
