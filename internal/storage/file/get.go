package file

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dsbasko/yandex-go-shortener/internal/entities"
)

// GetURLByOriginalURL returns a URL by original URL.
func (s *Storage) GetURLByOriginalURL(
	ctx context.Context,
	originalURL string,
) (resp entities.URL, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, err = s.file.Seek(0, 0)
	if err != nil {
		return entities.URL{}, fmt.Errorf("failed to seek to the beginning of the file: %w", err)
	}

	scanner := bufio.NewScanner(s.file)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return entities.URL{}, fmt.Errorf("failed to scan file: %w", scanner.Err())
		}

		select {
		case <-ctx.Done():
			return entities.URL{}, ctx.Err()
		default:
		}

		if !strings.Contains(scanner.Text(), `"original_url":"`+originalURL+`"`) {
			continue
		}

		if err = json.Unmarshal(scanner.Bytes(), &resp); err != nil {
			return entities.URL{}, fmt.Errorf("failed to unmarshal JSON data: %w", err)
		}
	}

	if resp.ShortURL == "" {
		return entities.URL{}, ErrURLNotFound
	}

	return resp, nil
}

// GetURLByShortURL returns a URL by short URL.
func (s *Storage) GetURLByShortURL(
	ctx context.Context,
	shortURL string,
) (resp entities.URL, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, err = s.file.Seek(0, 0)
	if err != nil {
		return entities.URL{}, fmt.Errorf("failed to seek to the beginning of the file: %w", err)
	}

	scanner := bufio.NewScanner(s.file)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return entities.URL{}, fmt.Errorf("failed to scan file: %w", scanner.Err())
		}

		select {
		case <-ctx.Done():
			return entities.URL{}, ctx.Err()
		default:
		}

		if !strings.Contains(scanner.Text(), `"short_url":"`+shortURL+`"`) {
			continue
		}

		if err = json.Unmarshal(scanner.Bytes(), &resp); err != nil {
			return entities.URL{}, fmt.Errorf("failed to unmarshal JSON data: %w", err)
		}
	}

	if resp.ShortURL == "" {
		return entities.URL{}, ErrURLNotFound
	}

	return resp, nil
}

// GetURLsByUserID returns URLs by user ID.
func (s *Storage) GetURLsByUserID(
	ctx context.Context,
	userID string,
) (resp []entities.URL, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, err = s.file.Seek(0, 0)
	if err != nil {
		return []entities.URL{}, fmt.Errorf("failed to seek to the beginning of the file: %w", err)
	}

	scanner := bufio.NewScanner(s.file)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return []entities.URL{}, fmt.Errorf("failed to scan file: %w", scanner.Err())
		}

		select {
		case <-ctx.Done():
			return []entities.URL{}, ctx.Err()
		default:
		}

		if !strings.Contains(scanner.Text(), `"user_id":"`+userID+`"`) {
			continue
		}

		found := entities.URL{}
		if err = json.Unmarshal(scanner.Bytes(), &found); err != nil {
			return []entities.URL{}, fmt.Errorf("failed to unmarshal JSON data: %w", err)
		}

		resp = append(resp, found)
	}

	if len(resp) == 0 {
		return []entities.URL{}, ErrURLsNotFound
	}

	return resp, nil
}
