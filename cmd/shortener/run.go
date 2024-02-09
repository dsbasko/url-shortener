package main

import (
	"context"
	"fmt"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/controller/rest"
	"github.com/dsbasko/yandex-go-shortener/internal/repository/storage"
	"github.com/dsbasko/yandex-go-shortener/internal/service/urls"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := config.Init()
	if err != nil {
		return fmt.Errorf("the configuration could not be loaded: %w", err)
	}

	log, err := logger.New(config.GetEnv(), "url-shortener")
	if err != nil {
		return fmt.Errorf("failed to load the logger: %w", err)
	}

	printGreeting(log)

	store, err := storage.New(ctx, log)
	if err != nil {
		return fmt.Errorf("storage could not be started: %w", err)
	}

	urlService := urls.New(log, store)

	if err = rest.New(ctx, log, store, urlService); err != nil {
		return fmt.Errorf("http server stopped with an error: %w", err)
	}

	return nil
}
