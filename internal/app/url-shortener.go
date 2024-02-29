package app

import (
	"context"

	"github.com/dsbasko/url-shortener/internal/config"
	httpController "github.com/dsbasko/url-shortener/internal/controller/http"
	storages "github.com/dsbasko/url-shortener/internal/repository/storage"
	"github.com/dsbasko/url-shortener/internal/service/urls"
	"github.com/dsbasko/url-shortener/pkg/graceful"
	"github.com/dsbasko/url-shortener/pkg/logger"
)

// RunURLShortener runs the REST server
func RunURLShortener(buildVersion, buildDate, buildCommit string) error {
	ctx, cancel := graceful.Context(context.Background(), graceful.DefaultSignals...)
	defer cancel()

	config.MustInit()
	log := logger.MustNew(config.Env(), "url-shortener")

	log.Infof("Build version: %s", buildVersion)
	log.Infof("Build date: %s", buildDate)
	log.Infof("Build commit: %s", buildCommit)

	storage := storages.MustNew(ctx, log)
	defer func() {
		if err := storage.Close(); err != nil {
			log.Errorf("storage could not be closed: %v", err)
		}
	}()

	urlService := urls.New(ctx, log, storage, storage, storage)
	httpController.New(ctx, log, storage, urlService)

	graceful.Wait()
	log.Infof("app has been stopped gracefully [%v]", graceful.Count())

	return nil
}
