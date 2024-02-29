package memory

import (
	"context"
	"fmt"
	"maps"

	"github.com/dsbasko/url-shortener/internal/entity"
)

// CreateURL creates a new URL.
func (s *Storage) CreateURL(
	_ context.Context,
	dto entity.URL,
) (resp entity.URL, unique bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if url, ok := s.storageOriginal[dto.OriginalURL]; ok {
		resp.OriginalURL = url.OriginalURL
		resp.ShortURL = url.ShortURL

		return resp, false, nil
	}

	s.storageShort[dto.ShortURL] = dto
	s.storageOriginal[dto.OriginalURL] = dto

	return dto, true, nil
}

// CreateURLs creates URLs.
func (s *Storage) CreateURLs(
	_ context.Context,
	dto []entity.URL,
) (resp []entity.URL, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	storageShortCopy := make(map[string]entity.URL, len(s.storageShort)+len(dto))
	maps.Copy(storageShortCopy, s.storageShort)

	storageOriginalCopy := make(map[string]entity.URL, len(s.storageOriginal)+len(dto))
	maps.Copy(storageOriginalCopy, s.storageOriginal)

	for _, url := range dto {
		if _, ok := storageOriginalCopy[url.OriginalURL]; ok {
			return []entity.URL{}, fmt.Errorf("url %s already exists", url.ShortURL)
		}

		storageShortCopy[url.ShortURL] = url
		storageOriginalCopy[url.OriginalURL] = url
	}

	s.storageShort = storageShortCopy
	s.storageOriginal = storageOriginalCopy

	return dto, nil
}
