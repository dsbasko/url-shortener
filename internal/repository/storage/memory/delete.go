package memory

import (
	"context"

	"github.com/dsbasko/yandex-go-shortener/internal/entity"
)

// DeleteURLs deletes URLs by user ID.
func (s *Storage) DeleteURLs(
	_ context.Context,
	dto []entity.URL,
) (resp []entity.URL, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, url := range dto {
		if foundURL, ok := s.storeShort[url.ShortURL]; ok && foundURL.UserID == url.UserID {
			url.OriginalURL = foundURL.OriginalURL
			foundURL.DeletedFlag = true
			s.storeShort[url.ShortURL] = foundURL
			resp = append(resp, foundURL)
		}

		if foundURL, ok := s.storeOriginal[url.OriginalURL]; ok && foundURL.UserID == url.UserID {
			foundURL.DeletedFlag = true
			s.storeOriginal[url.OriginalURL] = foundURL
		}
	}

	return resp, nil
}
