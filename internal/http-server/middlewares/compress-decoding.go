package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	middlewareChi "github.com/go-chi/chi/v5/middleware"
)

// CompressDecoding decompresses request.
func (m *Middlewares) CompressDecoding(next http.Handler) http.Handler {
	m.log.Debug("compress decoding middlewares enabled")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			m.log.Errorw(err.Error(), "request_id", middlewareChi.GetReqID(r.Context()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			if err = gz.Close(); err != nil {
				m.log.Debug(err.Error())
			}
		}()

		body, err := io.ReadAll(gz)
		if err != nil {
			m.log.Errorw(err.Error(), "request_id", middlewareChi.GetReqID(r.Context()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r)
	})
}
