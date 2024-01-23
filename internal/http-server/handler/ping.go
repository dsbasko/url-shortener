package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	log := h.log.With("request_id", middleware.GetReqID(r.Context()))

	if err := h.storage.Ping(r.Context()); err != nil {
		log.Errorw("no connection to the storage: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("pong")); err != nil {
		log.Errorf("failed to return response body: %v", err)
		return
	}
}
