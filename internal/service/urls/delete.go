package urls

import (
	"context"

	"github.com/dsbasko/yandex-go-shortener/internal/entity"
)

// URLDeleter is an interface for deleting URLs.
type URLDeleter interface {
	// DeleteURLs deletes URLs.
	DeleteURLs(ctx context.Context, dto []entity.URL) (resp []entity.URL, err error)
}

// DeleteURLs deletes urls from storage.
func (u *URLs) DeleteURLs(userID string, shortURLs []string) error {
	var urlDeleter URLDeleter = u.urlMutator

	urlsToDelete := make([]entity.URL, 0, len(shortURLs))
	for _, url := range shortURLs {
		urlsToDelete = append(urlsToDelete, entity.URL{
			ShortURL: url,
			UserID:   userID,
		})
	}

	go func() {
		if _, err := urlDeleter.DeleteURLs(context.Background(), urlsToDelete); err != nil {
			u.log.Errorw(err.Error())
			return
		}

		urls := make([]string, 0, len(urlsToDelete))
		for _, url := range urlsToDelete {
			urls = append(urls, url.ShortURL)
		}
		u.log.Infof("deleted urls: %v", urls)
	}()

	return nil
}
