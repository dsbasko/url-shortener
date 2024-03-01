package handlers

import "net/http"

// BadRequest returns a bad request.
func (h *Handler) BadRequest(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
