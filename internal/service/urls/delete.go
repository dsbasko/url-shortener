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
	u.deleteTask <- map[string][]string{
		userID: shortURLs,
	}

	return nil
}

// deleteWorker is a worker for deleting urls.
func (u *URLs) deleteWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			u.log.Infof("delete worker stopped by context")
			return
		case task := <-u.deleteTask:
			var urlDeleter URLDeleter = u.urlMutator

			urlsToDelete := make([]entity.URL, 0, len(task))
			for userID, shortURL := range task {
				for _, url := range shortURL {
					urlsToDelete = append(urlsToDelete, entity.URL{
						ShortURL: url,
						UserID:   userID,
					})
				}
			}

			deletedURLs, err := urlDeleter.DeleteURLs(context.Background(), urlsToDelete)
			if err != nil {
				u.log.Errorw(err.Error())
				return
			}

			urls := make([]string, 0, len(deletedURLs))
			for _, url := range deletedURLs {
				urls = append(urls, url.ShortURL)
			}

			if len(urls) > 0 {
				u.log.Debugf("deleted urls: %v", urls)
			}
		}
	}
}
