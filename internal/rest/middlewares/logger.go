package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	mwChi "github.com/go-chi/chi/v5/middleware"
)

// Logger sends request and response info to logger.
func (m *Middlewares) Logger(next http.Handler) http.Handler {
	m.log.Debug("logger middlewares enabled")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := m.log.With("request_id", mwChi.GetReqID(r.Context()))

		args := []any{
			"request_user_agent", r.UserAgent(),
			"request_id", mwChi.GetReqID(r.Context()),
			"request_method", r.Method,
			"request_path", r.URL.Path,
		}

		if body, err := io.ReadAll(r.Body); err == nil {
			args = append(args, "request_body", string(body))
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		buf := bytes.NewBuffer(nil)
		rw := &respWriter{w, buf}
		ww := mwChi.NewWrapResponseWriter(rw, r.ProtoMajor)
		timeStart := time.Now()

		defer func() {
			args = append(args, []any{
				"response_status", ww.Status(),
				"response_bytes", ww.BytesWritten(),
				"response_duration", time.Since(timeStart).String(),
			}...)

			if ww.Status() < http.StatusOK || ww.Status() >= http.StatusBadRequest || ww.BytesWritten() == 0 {
				args = append(args, "response_body", "null")
				m.log.Infow("request", args...)
				return
			}

			responseBuf := bytes.NewBuffer(nil)

			if !strings.Contains(ww.Header().Get("Content-Encoding"), "gzip") {
				if _, err := io.Copy(responseBuf, buf); err != nil {
					m.log.Infow("request", args...)
					return
				}

				args = append(args, "response_body", responseBuf.String())
				m.log.Infow("request", args...)
				return
			}

			reader, errGZIP := gzip.NewReader(buf)
			if errGZIP != nil {
				log.Debugf("failed to create gzip reader: %s", errGZIP.Error())
				m.log.Infow("request", args...)
				return
			}
			defer func() {
				if err := reader.Close(); err != nil {
					log.Debugf("failed to close gzip reader: %s", err.Error())
				}
			}()

			if _, err := io.Copy(responseBuf, reader); err != nil { //nolint:gosec
				log.Debugf("failed to copy gzip reader: %s", err.Error())
				m.log.Infow("request", args...)
				return
			}

			args = append(args, "response_body", responseBuf.String())
			m.log.Infow("request", args...)
		}()

		next.ServeHTTP(ww, r)
	})
}
