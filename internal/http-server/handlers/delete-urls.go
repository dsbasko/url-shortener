package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/dsbasko/yandex-go-shortener/internal/jwt"
)

// DeleteURLs deletes urls by user id.
func (h *Handler) DeleteURLs(w http.ResponseWriter, r *http.Request) {
	log := h.log.With("request_id", middleware.GetReqID(r.Context()))

	token, err := jwt.GetFromCookie(r)
	if err != nil || token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID := jwt.TokenToUserID(token)

	if r.ContentLength <= 4 { //nolint:gomnd
		h.log.Error(ErrEmptyBody)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var deleteURLs []string
	err = json.NewDecoder(r.Body).Decode(&deleteURLs)
	if err != nil {
		log.Errorw(fmt.Errorf("failed to unmarshal request body: %w", err).Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.urls.DeleteURLs(userID, deleteURLs)
	if err != nil {
		log.Errorw(fmt.Errorf("failed to delete urls: %w", err).Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
