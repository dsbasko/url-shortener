package memory

import (
	"context"
	"strconv"

	"github.com/dsbasko/yandex-go-shortener/internal/entity"
)

// Stats returns the stats of the URL.
func (s *Storage) Stats(_ context.Context) (entity.URLStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make(map[string]struct{})
	urlCount := 0

	for _, url := range s.storageShort {
		users[url.UserID] = struct{}{}
		urlCount++
	}

	return entity.URLStats{
		Users: strconv.Itoa(len(users)),
		URLs:  strconv.Itoa(urlCount),
	}, nil
}
