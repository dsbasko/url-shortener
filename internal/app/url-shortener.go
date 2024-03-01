package app

import (
	"context"
	"fmt"

	"github.com/dsbasko/url-shortener/internal/config"
	grpcController "github.com/dsbasko/url-shortener/internal/controller/grpc"
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

	switch config.Controller() {
	case "grpc":
		grpcController.Run(ctx, log, storage, urlService)
	case "http":
		httpController.Run(ctx, log, storage, urlService)
	default:
		return fmt.Errorf("unknown controller: %s", config.Controller())
	}

	graceful.Wait()
	return nil
}
