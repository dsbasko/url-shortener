package http

import (
	"context"
	"errors"
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	mwChi "github.com/go-chi/chi/v5/middleware"

	"github.com/dsbasko/url-shortener/internal/config"
	"github.com/dsbasko/url-shortener/internal/controller/http/handlers"
	"github.com/dsbasko/url-shortener/internal/controller/http/middlewares"
	"github.com/dsbasko/url-shortener/internal/service/urls"
	"github.com/dsbasko/url-shortener/pkg/graceful"
	"github.com/dsbasko/url-shortener/pkg/logger"
)

// Run creates a new http server.
func Run(
	ctx context.Context,
	log *logger.Logger,
	pinger handlers.Pinger,
	urlService urls.URLs,
) {
	router := chi.NewRouter()

	mw := middlewares.New(log)
	router.Use(mwChi.RequestID)
	router.Use(mwChi.Recoverer)
	router.Use(mwChi.RealIP)
	router.Use(mw.CompressDecoding)
	router.Use(mw.Logger)
	router.Use(mw.JWT)
	router.Use(mw.CompressEncoding)

	if config.IsEnabledPPROF() {
		router.Mount("/debug", mwChi.Profiler())
	}

	h := handlers.New(log, pinger, urlService)
	router.MethodNotAllowed(h.BadRequest)
	router.Get("/ping", h.Ping)
	router.Get("/{short_url}", h.Redirect)
	router.Get("/api/user/urls", h.GetURLsByUserID)
	router.Post("/", h.CreateURLTextPlain)
	router.Post("/api/shorten", h.CreateURLJSON)
	router.Post("/api/shorten/batch", h.CreateURLsJSON)
	router.Delete("/api/user/urls", h.DeleteURLs)
	router.Get("/api/internal/stats", h.Stats)

	routes := router.Routes()
	for _, route := range routes {
		for handle := range route.Handlers {
			log.Debugf("mapped [%v] %v route", handle, route.Pattern)
		}
	}

	server := http.Server{
		Addr:         config.ServerAddress(),
		Handler:      router,
		ReadTimeout:  config.RESTReadTimeout(),
		WriteTimeout: config.RESTWriteTimeout(),
	}

	// Run the server in a goroutine so that it doesn't block.
	// The server will be gracefully shutdown by the signal.
	graceful.Add()
	go runServer(log, &server)

	graceful.Add()
	go gracefulShutdown(ctx, log, &server)
}

// runServer runs the http server.
func runServer(log *logger.Logger, server *http.Server) {
	defer graceful.Done()

	if config.IsRESTEnableHTTPS() {
		log.Infof("starting rest server at the address: https://%s", config.ServerAddress())
		err := server.ListenAndServeTLS(path.Join("cert", "cert.pem"), path.Join("cert", "key.pem"))
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("failed to start the rest server: %v", err)
		}
		return
	}

	log.Infof("starting rest server at the address: http://%s", config.ServerAddress())
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Errorf("failed to start the rest server: %v", err)
	}
}

// gracefulShutdown gracefully shutdowns the http server.
func gracefulShutdown(ctx context.Context, log *logger.Logger, server *http.Server) {
	defer graceful.Done()

	<-ctx.Done()

	server.SetKeepAlivesEnabled(false)
	log.Infof("shutdown http server by signal")

	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("a signal has been received to terminate the http server: %v", err)
	}
}
