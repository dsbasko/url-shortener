package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dsbasko/yandex-go-shortener/pkg/api"
	"github.com/go-chi/chi/v5/middleware"
)

// CreateURLsJSON creates url with json body.
func (h *Handler) CreateURLJSON(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var dto api.CreateURLRequest

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

	createdURL, unique, err := h.urls.CreateURL(r.Context(), dto.URL)
	if err != nil {
		log.Errorw(fmt.Errorf("failed to create link in urls: %w", err).Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := api.CreateURLResponse{
		Result: createdURL.ShortURL,
	}

	responseBytes, responseErr := json.Marshal(response)
	if responseErr != nil {
		log.Errorw(fmt.Errorf("failed to assemble the structure into json: %w", responseErr).Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if unique {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusConflict)
	}

	if _, err = w.Write(responseBytes); err != nil {
		log.Errorw(fmt.Errorf("failed to return response body: %w", err).Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
