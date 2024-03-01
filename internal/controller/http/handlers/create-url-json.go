package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	api "github.com/dsbasko/url-shortener/api/http"
)

// CreateURLJSON creates url with json body.
func (h *Handler) CreateURLJSON(w http.ResponseWriter, r *http.Request) {
	var dto api.CreateURLRequest

	log := h.log.With("request_id", middleware.GetReqID(r.Context()))

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		log.Errorf("failed to decode json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdURL, unique, err := h.urls.CreateURL(r.Context(), dto.URL)
	if err != nil {
		log.Errorf("failed to create link in urls: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := api.CreateURLResponse{
		Result: createdURL.ShortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	if unique {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusConflict)
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Errorf("failed to return response body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
