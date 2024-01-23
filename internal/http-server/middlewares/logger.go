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

type responseLogger struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (r *responseLogger) Write(b []byte) (int, error) {
	r.buf.Write(b)
	return r.ResponseWriter.Write(b)
}

func (m *Middlewares) Logger(next http.Handler) http.Handler {
	m.log.Debug("logger middlewares enabled")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		args := []any{
			"request_method", r.Method,
			"request_path", r.URL.Path,
			"request_user_agent", r.UserAgent(),
			"request_id", middlewareChi.GetReqID(r.Context()),
		}

		if r.Body != nil {
			body, err := io.ReadAll(r.Body)
			if err == nil {
				args = append(args, "request_body", string(body))
				r.Body = io.NopCloser(bytes.NewBuffer(body))
			}
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
					reader, err := gzip.NewReader(buf)
					if err == nil {
						defer reader.Close()
						if _, err = io.Copy(responseBuf, reader); err != nil { //nolint:gosec
							next.ServeHTTP(ww, r)
						}
					}
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
