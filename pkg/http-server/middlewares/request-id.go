package middlewares

import (
	"net/http"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

// RequestID adds request id to header.
func (m *Middlewares) RequestID(next http.Handler) http.Handler {
	m.log.Debug("Adding request id to header is enabled")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("request-id", chiMiddleware.GetReqID(r.Context()))
		next.ServeHTTP(w, r)
	})
}
