package file

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/dsbasko/url-shortener/internal/entity"
)

// Stats returns the stats of the URL.
func (s *Storage) Stats(ctx context.Context) (entity.URLStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make(map[string]struct{})
	urlCount := 0

	_, err := s.file.Seek(0, 0)
	if err != nil {
		return entity.URLStats{}, fmt.Errorf("failed to seek to the beginning of the file: %w", err)
	}

	scanner := bufio.NewScanner(s.file)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return entity.URLStats{}, fmt.Errorf("failed to scan file: %w", scanner.Err())
		}

		if ctx.Err() != nil {
			return entity.URLStats{}, ctx.Err()
		}

		url := entity.URL{}
		if err = json.Unmarshal(scanner.Bytes(), &url); err != nil {
			return entity.URLStats{}, fmt.Errorf("failed to unmarshal JSON data: %w", err)
		}

		users[url.UserID] = struct{}{}
		urlCount++
	}

	return entity.URLStats{
		Users: strconv.Itoa(len(users)),
		URLs:  strconv.Itoa(urlCount),
	}, nil
}
