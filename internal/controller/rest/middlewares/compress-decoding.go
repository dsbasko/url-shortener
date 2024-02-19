package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	mwChi "github.com/go-chi/chi/v5/middleware"
)

// CompressDecoding decompresses request.
func (m *Middlewares) CompressDecoding(next http.Handler) http.Handler {
	m.log.Debug("compress decoding middlewares enabled")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := m.log.With("request_id", mwChi.GetReqID(r.Context()))

		if r.ContentLength == 0 {
			next.ServeHTTP(w, r)
			return
		}

		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			log.Debugf("failed to create gzip reader: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer func() {
			if err = gz.Close(); err != nil {
				log.Debugf("failed to close gzip reader: %s", err.Error())
			}
		}()

		r.Body = io.NopCloser(gz)
		next.ServeHTTP(w, r)
	})
}
