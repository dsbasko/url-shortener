package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/dsbasko/url-shortener/internal/service/jwt"
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

	var deleteURLs []string
	err = json.NewDecoder(r.Body).Decode(&deleteURLs)
	if err != nil {
		log.Errorf("failed to unmarshal request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.urls.DeleteURLs(userID, deleteURLs)
	if err != nil {
		log.Errorf("failed to delete urls: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
