package memory

import (
	"context"
	"fmt"
	"maps"

	"github.com/dsbasko/yandex-go-shortener/internal/entity"
)

// CreateURL creates a new URL.
func (s *Storage) CreateURL(
	_ context.Context,
	dto entity.URL,
) (resp entity.URL, unique bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if url, ok := s.storeOriginal[dto.OriginalURL]; ok {
		resp.OriginalURL = url.OriginalURL
		resp.ShortURL = url.ShortURL

		return resp, false, nil
	}

	s.storeShort[dto.ShortURL] = dto
	s.storeOriginal[dto.OriginalURL] = dto

	return dto, true, nil
}

// CreateURLs creates URLs.
func (s *Storage) CreateURLs(
	_ context.Context,
	dto []entity.URL,
) (resp []entity.URL, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	storeShortCopy := make(map[string]entity.URL, len(s.storeShort)+len(dto))
	maps.Copy(storeShortCopy, s.storeShort)

	storeOriginalCopy := make(map[string]entity.URL, len(s.storeOriginal)+len(dto))
	maps.Copy(storeOriginalCopy, s.storeOriginal)

	for _, url := range dto {
		if _, ok := storeOriginalCopy[url.OriginalURL]; ok {
			return []entity.URL{}, fmt.Errorf("url %s already exists", url.ShortURL)
		}

		storeShortCopy[url.ShortURL] = url
		storeOriginalCopy[url.OriginalURL] = url
	}

	s.storeShort = storeShortCopy
	s.storeOriginal = storeOriginalCopy

	return dto, nil
}
