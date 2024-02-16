package app

import (
	"context"
	"fmt"

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

	err := config.Init()
	if err != nil {
		return fmt.Errorf("the configuration could not be loaded: %w", err)
	}

	log, err := logger.New(config.Env(), "url-shortener")
	if err != nil {
		return fmt.Errorf("failed to load the logger: %w", err)
	}

	log.Infof("Build version: %s", buildVersion)
	log.Infof("Build date: %s", buildDate)
	log.Infof("Build commit: %s", buildCommit)

	storage, err := storages.New(ctx, log)
	if err != nil {
		return fmt.Errorf("storage could not be started: %w", err)
	}
	defer func() {
		if err = storage.Close(); err != nil {
			log.Errorf("storage could not be closed: %v", err)
		}
	}()

	urlService := urls.New(ctx, log, storage, storage)
	rest.New(ctx, log, storage, urlService)

	graceful.Wait()
	log.Infof("app has been stopped gracefully [%v]", graceful.Count())

	return nil
}
