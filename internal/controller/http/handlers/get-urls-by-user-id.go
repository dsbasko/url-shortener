package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/dsbasko/url-shortener/internal/service/jwt"
)

// GetURLsByUserID returns all urls by user id.
func (h *Handler) GetURLsByUserID(w http.ResponseWriter, r *http.Request) {
	log := h.log.With("request_id", middleware.GetReqID(r.Context()))

	token, err := jwt.GetFromCookie(r)
	if err != nil || token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID := jwt.TokenToUserID(token)

	urlsResp, err := h.urls.GetURLsByUserID(r.Context(), userID)
	if err != nil {
		log.Errorf("failed to get link from urls: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(urlsResp) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(urlsResp); err != nil {
		log.Errorf("failed to return response body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
