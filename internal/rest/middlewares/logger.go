package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	middlewareChi "github.com/go-chi/chi/v5/middleware"
)

// responseLogger is a response logger.
type responseLogger struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

// Write writes response.
func (r *responseLogger) Write(b []byte) (int, error) {
	r.buf.Write(b)
	return r.ResponseWriter.Write(b)
}

// Logger sends request and response info to logger.
func (m *Middlewares) Logger(next http.Handler) http.Handler {
	m.log.Debug("logger middlewares enabled")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		args := []any{
			"request_method", r.Method,
			"request_path", r.URL.Path,
			"request_user_agent", r.UserAgent(),
			"request_id", middlewareChi.GetReqID(r.Context()),
		}

		if body, err := io.ReadAll(r.Body); err == nil {
			args = append(args, "request_body", string(body))
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		buf := bytes.NewBuffer(nil)
		rl := &responseLogger{w, buf}
		ww := middlewareChi.NewWrapResponseWriter(rl, r.ProtoMajor)
		timeStart := time.Now()

		defer func() {
			args = append(args, []any{
				"response_status", ww.Status(),
				"response_bytes", ww.BytesWritten(),
				"response_duration", time.Since(timeStart).String(),
			}...)

			if ww.Status() >= http.StatusOK && ww.Status() < http.StatusBadRequest {
				responseBuf := bytes.NewBuffer(nil)

				if strings.Contains(ww.Header().Get("Content-Encoding"), "gzip") {
					reader, errGZIP := gzip.NewReader(buf)
					if errGZIP != nil {
						m.log.Warnw("failed to create gzip reader", "error", errGZIP)
						next.ServeHTTP(ww, r)
					}

					defer func() {
						if err := reader.Close(); err != nil {
							m.log.Errorw("failed to close gzip reader", "error", err)
						}
					}()

					if _, err := io.Copy(responseBuf, reader); err != nil { //nolint:gosec
						m.log.Warnw("failed to copy gzip reader", "error", err)
						next.ServeHTTP(ww, r)
					}

					defer func(reader *gzip.Reader) {
						if err := reader.Close(); err != nil {
							m.log.Errorw("failed to close gzip reader", "error", err)
						}
					}(reader)
				} else {
					if _, err := io.Copy(responseBuf, buf); err != nil {
						next.ServeHTTP(ww, r)
					}
				}

				args = append(args, "response_body", responseBuf.String())
			} else {
				args = append(args, "response_body", "null")
			}

			m.log.Infow("request", args...)
		}()

		next.ServeHTTP(ww, r)
	})
}
