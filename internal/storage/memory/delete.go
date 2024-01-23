package memory

import (
	"context"

	"github.com/dsbasko/yandex-go-shortener/internal/entities"
)

func (s *Storage) DeleteURLs(
	ctx context.Context,
	dto []entities.URL,
) (resp []entities.URL, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, url := range dto {
		select {
		case <-ctx.Done():
			return []entities.URL{}, ctx.Err()
		default:
		}

		if foundURL, ok := s.store[url.ShortURL]; ok && foundURL.UserID == url.UserID {
			foundURL.DeletedFlag = true
			s.store[url.ShortURL] = foundURL
			resp = append(resp, foundURL)
		}
	}

	return resp, nil
}
