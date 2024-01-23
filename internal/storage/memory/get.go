package memory

import (
	"context"

	"github.com/dsbasko/yandex-go-shortener/internal/entities"
)

func (s *Storage) GetURLByOriginalURL(
	_ context.Context,
	originalURL string,
) (resp entities.URL, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, url := range s.store {
		if url.OriginalURL == originalURL {
			return url, nil
		}
	}

	return entities.URL{}, nil
}

func (s *Storage) GetURLByShortURL(
	_ context.Context,
	shortURL string,
) (resp entities.URL, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if url, ok := s.store[shortURL]; ok {
		return url, nil
	}

	return entities.URL{}, nil
}

func (s *Storage) GetURLsByUserID(
	_ context.Context,
	userID string,
) (resp []entities.URL, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, url := range s.store {
		if url.UserID == userID {
			resp = append(resp, url)
		}
	}

	return resp, nil
}