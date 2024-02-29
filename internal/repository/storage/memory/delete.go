package memory

import (
	"context"

	"github.com/dsbasko/url-shortener/internal/entity"
)

// DeleteURLs deletes URLs by user ID.
func (s *Storage) DeleteURLs(
	_ context.Context,
	dto []entity.URL,
) (resp []entity.URL, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, url := range dto {
		if foundURL, ok := s.storageShort[url.ShortURL]; ok && foundURL.UserID == url.UserID {
			url.OriginalURL = foundURL.OriginalURL
			foundURL.DeletedFlag = true
			s.storageShort[url.ShortURL] = foundURL
			resp = append(resp, foundURL)
		}

		if foundURL, ok := s.storageOriginal[url.OriginalURL]; ok && foundURL.UserID == url.UserID {
			foundURL.DeletedFlag = true
			s.storageOriginal[url.OriginalURL] = foundURL
		}
	}

	return resp, nil
}
