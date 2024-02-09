package rest

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	mwChi "github.com/go-chi/chi/v5/middleware"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/controller/rest/handlers"
	"github.com/dsbasko/yandex-go-shortener/internal/controller/rest/middlewares"
	"github.com/dsbasko/yandex-go-shortener/internal/service/urls"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

// New creates a new http server.
func New(
	ctx context.Context,
	log *logger.Logger,
	pinger handlers.Pinger,
	urlService urls.URLs,
) error {
	router := chi.NewRouter()

	mw := middlewares.New(log)
	router.Use(mwChi.RequestID)
	router.Use(mwChi.Recoverer)
	router.Use(mw.CompressDecoding)
	router.Use(mw.Logger)
	router.Use(mw.JWT)
	router.Use(mw.CompressEncoding)

	router.Mount("/debug", mwChi.Profiler())

	h := handlers.New(log, pinger, urlService)
	router.Get("/ping", h.Ping)
	router.Get("/{short_url}", h.Redirect)
	router.Get("/api/user/urls", h.GetURLsByUserID)
	router.Post("/", h.CreateURLTextPlain)
	router.Post("/api/shorten", h.CreateURLJSON)
	router.Post("/api/shorten/batch", h.CreateURLsJSON)
	router.Delete("/api/user/urls", h.DeleteURLs)

	routes := router.Routes()
	for _, route := range routes {
		for handle := range route.Handlers {
			log.Debugf("mapped [%v] %v route", handle, route.Pattern)
		}
	}

	server := http.Server{
		Addr:         config.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  config.GetRestReadTimeout(),
		WriteTimeout: config.GetRestWriteTimeout(),
	}

	go func() {
		<-ctx.Done()
		log.Info("shutdown rest server")
		server.SetKeepAlivesEnabled(false)
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Errorf("a signal has been received to terminate the http server: %v", err)
		}
	}()

	log.Infof("starting rest server at the address: %s", config.GetServerAddress())
	return server.ListenAndServe()
}
