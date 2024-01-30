package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// CreateURLTextPlain creates url with text/plain body.
func (h *Handler) CreateURLTextPlain(w http.ResponseWriter, r *http.Request) {
	log := h.log.With("request_id", middleware.GetReqID(r.Context()))

	if r.ContentLength <= 4 { //nolint:gomnd
		h.log.Error(ErrEmptyBody)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalURL, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorw(fmt.Errorf("failed to read request body: %w", err).Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdURL, unique, err := h.urls.CreateURL(r.Context(), string(originalURL))
	if err != nil {
		log.Errorw(fmt.Errorf("failed to create link in urls: %w", err).Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if unique {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusConflict)
	}

	if _, err = w.Write([]byte(createdURL.ShortURL)); err != nil {
		log.Errorw(fmt.Errorf("failed to return response body: %w", err).Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
