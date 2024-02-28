package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
)

// Stats returns the stats of the URL.
func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	log := h.log.With("request_id", middleware.GetReqID(r.Context()))

	isTrustedSubnet, err := config.IsTrustedSubnet(r.RemoteAddr)
	if err != nil {
		log.Errorf("failed to check trusted subnet: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !isTrustedSubnet {
		log.Warnf("untrusted subnet: %s", r.RemoteAddr)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	stats, err := h.urls.Stats(r.Context())
	if err != nil {
		log.Errorf("failed to get stats from service layer: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(stats); err != nil {
		log.Errorf("failed to return response body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
