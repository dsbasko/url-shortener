package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// CreateURLTextPlain creates url with text/plain body.
func (h *Handler) CreateURLTextPlain(w http.ResponseWriter, r *http.Request) {
	log := h.log.With("request_id", middleware.GetReqID(r.Context()))

	originalURL, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("failed to read request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdURL, unique, err := h.urls.CreateURL(r.Context(), string(originalURL))
	if err != nil {
		log.Errorf("failed to create link in urls: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if unique {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusConflict)
	}
	_, _ = w.Write([]byte(createdURL.ShortURL))
}
