package middlewares

import (
	"bytes"
	"net/http"

	"github.com/dsbasko/url-shortener/pkg/logger"
)

// Middlewares a collection of middlewares.
type Middlewares struct {
	log *logger.Logger
}

// New creates a new middlewares constructor.
func New(log *logger.Logger) *Middlewares {
	return &Middlewares{
		log: log,
	}
}

// respWriter is a response to compress encoder.
type respWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

// respWriter writes response.
func (r *respWriter) Write(b []byte) (int, error) {
	r.buf.Write(b)
	return r.ResponseWriter.Write(b)
}
