package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// Pinger checks connection to storage.
type Pinger interface {
	// Ping checks connection to storage.
	Ping(ctx context.Context) error
}

// Ping returns pong if the server and storage are available.
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	log := h.log.With("request_id", middleware.GetReqID(r.Context()))

	if err := h.pinger.Ping(r.Context()); err != nil {
		log.Errorw("no connection to the storage: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("pong")); err != nil {
		log.Errorf("failed to return response body: %v", err)
		return
	}
}

// Generate mocks for tests.
//go:generate ../../../../bin/mockgen -destination=./mocks/pinger.go -package=mock_handlers github.com/dsbasko/yandex-go-shortener/internal/controller/rest/handlers Pinger
