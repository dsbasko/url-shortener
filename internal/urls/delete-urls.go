package urls

import (
	"context"

	"github.com/dsbasko/yandex-go-shortener/internal/entities"
)

// DeleteURLs deletes urls from storage.
func (u *URLs) DeleteURLs(userID string, shortURLs []string) error {
	urlsToDelete := make([]entities.URL, 0, len(shortURLs))
	for _, url := range shortURLs {
		urlsToDelete = append(urlsToDelete, entities.URL{
			ShortURL: url,
			UserID:   userID,
		})
	}

	go func() {
		if _, err := u.storage.DeleteURLs(context.Background(), urlsToDelete); err != nil {
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
