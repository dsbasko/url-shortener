package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	middlewareChi "github.com/go-chi/chi/v5/middleware"
)

// compressGzipWriter is a response logger.
type compressGzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write writes response.
func (w compressGzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// CompressEncoding compresses response.
func (m *Middlewares) CompressEncoding(next http.Handler) http.Handler {
	m.log.Debug("compress encoding middlewares enabled")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzWriter, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			m.log.Errorw(err.Error(), "request_id", middlewareChi.GetReqID(r.Context()))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer func() {
			if err = gzWriter.Close(); err != nil {
				m.log.Debug(err.Error())
			}
		}()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(compressGzipWriter{ResponseWriter: w, Writer: gzWriter}, r)
	})
}
