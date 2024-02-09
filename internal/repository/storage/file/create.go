package file

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/dsbasko/yandex-go-shortener/internal/entity"
)

// CreateURL creates a new URL.
func (s *Storage) CreateURL(
	ctx context.Context,
	dto entity.URL,
) (resp entity.URL, unique bool, err error) {
	if foundURL, errURL := s.GetURLByOriginalURL(ctx, dto.OriginalURL); errURL == nil && len(foundURL.OriginalURL) > 0 {
		return foundURL, false, nil
	}

	newID := s.getLastID()
	data, err := json.Marshal(entity.URL{
		ID:          newID,
		ShortURL:    dto.ShortURL,
		OriginalURL: dto.OriginalURL,
		UserID:      dto.UserID,
	})
	if err != nil {
		return entity.URL{}, false, fmt.Errorf("failed to marshal URL to JSON: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err = s.writer.Write(append(data, '\n')); err != nil {
		return entity.URL{}, false, fmt.Errorf("failed to write data to file: %w", err)
	}

	if err = s.writer.Flush(); err != nil {
		return entity.URL{}, false, fmt.Errorf("failed to flush writer: %w", err)
	}

	resp.ID = newID
	resp.ShortURL = dto.ShortURL
	resp.OriginalURL = dto.OriginalURL
	resp.UserID = dto.UserID

	return resp, true, nil
}

// CreateURLs creates URLs.
func (s *Storage) CreateURLs(
	ctx context.Context,
	dto []entity.URL,
) (resp []entity.URL, err error) {
	s.mu.Lock()

	if _, err = s.file.Seek(0, 0); err != nil {
		return []entity.URL{}, fmt.Errorf("failed to seek file: %w", err)
	}

	data, err := io.ReadAll(s.file)
	if err != nil {
		return []entity.URL{}, fmt.Errorf("failed to read file: %w", err)
	}

	backup := data
	s.mu.Unlock()

	lastID, err := strconv.Atoi(s.getLastID())
	if err != nil {
		return []entity.URL{}, fmt.Errorf("failed to convert string to int: %w", err)
	}

	for _, url := range dto {
		lastIDString := strconv.Itoa(lastID)

		foundURL, errFor := s.GetURLByOriginalURL(ctx, url.OriginalURL)
		if errFor == nil && foundURL.ShortURL != "" {
			return []entity.URL{}, s.rollbackCreateURLMany(
				backup,
				fmt.Errorf("%s: %s", ErrURLFoundDuplicate, url.OriginalURL),
			)
		}

		s.mu.Lock()
		data, err = json.Marshal(entity.URL{
			ID:          lastIDString,
			ShortURL:    url.ShortURL,
			OriginalURL: url.OriginalURL,
			UserID:      url.UserID,
		})
		if err != nil {
			return []entity.URL{}, s.rollbackCreateURLMany(
				backup,
				fmt.Errorf("failed to marshal URL to JSON: %w", err),
			)
		}

		if _, err = s.writer.Write(append(data, '\n')); err != nil {
			return []entity.URL{}, s.rollbackCreateURLMany(
				backup,
				fmt.Errorf("failed to write data to file: %w", err),
			)
		}

		s.mu.Unlock()
		resp = append(resp, entity.URL{
			ID:          lastIDString,
			ShortURL:    url.ShortURL,
			OriginalURL: url.OriginalURL,
			UserID:      url.UserID,
		})

		lastID++
	}

	if err = s.writer.Flush(); err != nil {
		return []entity.URL{}, s.rollbackCreateURLMany(
			backup,
			fmt.Errorf(": %w", err),
		)
	}

	return resp, nil
}

func (s *Storage) rollbackCreateURLMany(backup []byte, origErr error) error {
	if err := s.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	if err := s.file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}

	if _, err := s.file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	if _, err := s.writer.Write(backup); err != nil {
		return fmt.Errorf("failed to write backup data: %w", err)
	}

	if err := s.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return origErr
}
