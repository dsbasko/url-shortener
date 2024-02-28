package app

import (
	"context"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/controller/rest"
	storages "github.com/dsbasko/yandex-go-shortener/internal/repository/storage"
	"github.com/dsbasko/yandex-go-shortener/internal/service/urls"
	"github.com/dsbasko/yandex-go-shortener/pkg/graceful"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

// RunREST runs the REST server
func RunREST(buildVersion, buildDate, buildCommit string) error {
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
	rest.New(ctx, log, storage, urlService)

	graceful.Wait()
	log.Infof("app has been stopped gracefully [%v]", graceful.Count())

	return nil
}
