package memory

import (
	"context"

	"github.com/dsbasko/url-shortener/internal/entity"
)

// GetURLByOriginalURL returns a URL by original URL.
func (s *Storage) GetURLByOriginalURL(
	_ context.Context,
	originalURL string,
) (resp entity.URL, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if url, ok := s.storageOriginal[originalURL]; ok {
		return url, nil
	}
	return entity.URL{}, nil
}

// GetURLByShortURL returns a URL by short URL.
func (s *Storage) GetURLByShortURL(
	_ context.Context,
	shortURL string,
) (resp entity.URL, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if url, ok := s.storageShort[shortURL]; ok {
		return url, nil
	}

	return entity.URL{}, nil
}

// GetURLsByUserID returns URLs by user ID.
func (s *Storage) GetURLsByUserID(
	_ context.Context,
	userID string,
) (resp []entity.URL, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, url := range s.storageShort {
		if url.UserID == userID {
			resp = append(resp, url)
		}
	}

	return resp, nil
}
