package urls

import (
	"context"
	"fmt"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/entities"
)

func (u *URLs) GetURLsByUserID(ctx context.Context, userID string) ([]entities.URL, error) {
	storeResp, err := u.storage.GetURLsByUserID(ctx, userID)
	if err != nil {
		return []entities.URL{}, fmt.Errorf("error getting url from storage: %w", err)
	}

	var resp []entities.URL
	for _, url := range storeResp {
		url.ShortURL = fmt.Sprintf("%s%s", config.GetBaseURL(), url.ShortURL)
		resp = append(resp, url)
	}

	return resp, nil
}
