package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dsbasko/yandex-go-shortener/internal/jwt"
	"github.com/go-chi/chi/v5/middleware"
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

	urlsBytes, err := json.Marshal(urlsResp)
	if err != nil {
		log.Errorf("failed to marshal urls: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(urlsBytes); err != nil {
		log.Errorf("failed to write response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
