package file

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dsbasko/yandex-go-shortener/internal/entity"
)

// DeleteURLs returns a URL by original URL.
func (s *Storage) DeleteURLs(
	ctx context.Context,
	dto []entity.URL,
) (resp []entity.URL, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err = s.file.Seek(0, 0); err != nil {
		return []entity.URL{}, fmt.Errorf("failed to seek to the beginning of the file: %w", err)
	}

	var data []string
	var found bool
	scanner := bufio.NewScanner(s.file)

	for scanner.Scan() {
		if scanner.Err() != nil {
			return []entity.URL{}, fmt.Errorf("failed to scan file: %w", scanner.Err())
		}

		if ctx.Err() != nil {
			return []entity.URL{}, ctx.Err()
		}

		dataText := scanner.Text()
		for _, url := range dto {
			if !strings.Contains(dataText, fmt.Sprintf(
				`"user_id":%q,"short_url":%q`,
				url.UserID, url.ShortURL,
			)) {
				continue
			}

			var dataJSON entity.URL
			dataText = strings.ReplaceAll(dataText, `"is_deleted":false`, `"is_deleted":true`)
			found = true

			if err = json.Unmarshal([]byte(dataText), &dataJSON); err != nil {
				return []entity.URL{}, fmt.Errorf("failed to unmarshal JSON data: %w", err)
			}

			resp = append(resp, dataJSON)
		}

		data = append(data, dataText)
	}

	if !found {
		return []entity.URL{}, nil
	}

	if err = s.file.Truncate(0); err != nil {
		return []entity.URL{}, fmt.Errorf("failed to truncate file: %w", err)
	}

	if _, err = s.file.Seek(0, 0); err != nil {
		return []entity.URL{}, fmt.Errorf("failed to seek to the beginning of the file: %w", err)
	}

	for _, line := range data {
		if _, err = s.writer.Write(append([]byte(line), '\n')); err != nil {
			return []entity.URL{}, fmt.Errorf("failed to write data to file: %w", err)
		}
	}

	if err = s.writer.Flush(); err != nil {
		return []entity.URL{}, fmt.Errorf("failed to flush writer: %w", err)
	}

	return resp, nil
}
