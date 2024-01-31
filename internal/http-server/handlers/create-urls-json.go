package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/dsbasko/yandex-go-shortener/pkg/api"
)

// CreateURLsJSON creates urls with json body.
func (h *Handler) CreateURLsJSON(w http.ResponseWriter, r *http.Request) {
	var dto []api.CreateURLsRequest

	log := h.log.With("request_id", middleware.GetReqID(r.Context()))

	if r.Header.Get("Content-Type") != "application/json" {
		h.log.Error(ErrWrongContentType)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.ContentLength <= 4 { //nolint:gomnd
		h.log.Error(ErrEmptyBody)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		log.Errorw(fmt.Errorf("failed to decode json: %w", err).Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdURLs, err := h.urls.CreateURLs(r.Context(), dto)
	if err != nil {
		log.Errorw(fmt.Errorf("failed to create link in urls: %w", err).Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(createdURLs); err != nil {
		log.Errorw(fmt.Errorf("failed to return response body: %w", err).Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
