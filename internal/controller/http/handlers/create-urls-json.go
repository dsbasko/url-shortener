package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	api "github.com/dsbasko/url-shortener/api/http"
)

// CreateURLsJSON creates urls with json body.
func (h *Handler) CreateURLsJSON(w http.ResponseWriter, r *http.Request) {
	var dto []api.CreateURLsRequest

	log := h.log.With("request_id", middleware.GetReqID(r.Context()))

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		log.Errorf("failed to decode json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdURLs, err := h.urls.CreateURLs(r.Context(), dto)
	if err != nil {
		log.Errorf("failed to create link in urls: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(createdURLs); err != nil {
		log.Errorf("failed to return response body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
