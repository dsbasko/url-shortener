package memory

import (
	"context"
	"fmt"
	"maps"

	"github.com/dsbasko/yandex-go-shortener/internal/entities"
)

func (s *Storage) CreateURL(
	ctx context.Context,
	dto entities.URL,
) (resp entities.URL, unique bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if url, ok := s.store[dto.ShortURL]; ok {
		select {
		case <-ctx.Done():
			return resp, false, ctx.Err()
		default:
		}

		resp.OriginalURL = url.OriginalURL
		resp.ShortURL = dto.ShortURL
		return resp, false, nil
	}

	s.store[dto.ShortURL] = dto

	return dto, true, nil
}

func (s *Storage) CreateURLs(
	ctx context.Context,
	dto []entities.URL,
) (resp []entities.URL, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	storeCopy := make(map[string]entities.URL, len(s.store)+len(dto))
	maps.Copy(storeCopy, s.store)

	for _, url := range dto {
		select {
		case <-ctx.Done():
			return resp, ctx.Err()
		default:
		}

		if _, ok := storeCopy[url.ShortURL]; ok {
			return []entities.URL{}, fmt.Errorf("url %s already exists", url.ShortURL)
		}

		storeCopy[url.ShortURL] = url
	}

	s.store = storeCopy
	return dto, nil
}
