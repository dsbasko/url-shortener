package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	mwChi "github.com/go-chi/chi/v5/middleware"
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
		log := m.log.With("request_id", mwChi.GetReqID(r.Context()))

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		buf := bytes.NewBuffer(nil)
		rw := &respWriter{w, buf}
		ww := mwChi.NewWrapResponseWriter(rw, r.ProtoMajor)

		gzWriter, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			log.Debugf("failed to create gzip writer: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			next.ServeHTTP(w, r)
			return
		}
		defer func(gzWriter *gzip.Writer) {
			if err = gzWriter.Close(); err != nil && ww.BytesWritten() != 0 {
				log.Debugf("failed to close gzip writer: %s", err.Error())
			}
		}(gzWriter)

		w.Header().Set("Content-Encoding", "gzip")
		compressedWriter := compressGzipWriter{ResponseWriter: w, Writer: gzWriter}
		next.ServeHTTP(compressedWriter, r)
	})
}
