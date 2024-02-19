package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Redirect redirects to original url.
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	log := h.log.With("request_id", middleware.GetReqID(r.Context()))
	shortURL := chi.URLParam(r, "short_url")

	urlResp, err := h.urls.GetURL(r.Context(), shortURL)
	if err != nil {
		log.Errorf("failed to get link from urls: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if urlResp.OriginalURL == "" {
		log.Warnf("original url not found: %s", shortURL)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if urlResp.DeletedFlag {
		w.WriteHeader(http.StatusGone)
		return
	}

	http.Redirect(w, r, urlResp.OriginalURL, http.StatusTemporaryRedirect)
}
