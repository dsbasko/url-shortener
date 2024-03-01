package urls

import (
	"context"
	"fmt"

	"github.com/dsbasko/url-shortener/internal/entity"
)

// Analyzer represents the metric service.
type Analyzer interface {
	// Stat returns the stats of the URL.
	Stats(ctx context.Context) (resp entity.URLStats, err error)
}

// Stats returns the stats of the URL.
func (u *URLs) Stats(ctx context.Context) (entity.URLStats, error) {
	stats, err := u.urlAnalyzer.Stats(ctx)
	if err != nil {
		return entity.URLStats{}, fmt.Errorf("failed to get stats from repository layer: %w", err)
	}

	return stats, nil
}
