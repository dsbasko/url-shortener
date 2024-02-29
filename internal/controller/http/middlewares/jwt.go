package middlewares

import (
	"context"
	"errors"
	"net/http"

	mwChi "github.com/go-chi/chi/v5/middleware"

	"github.com/dsbasko/url-shortener/internal/service/jwt"
	"github.com/dsbasko/url-shortener/pkg/logger"
)

// JWT adds jwt token to request context.
func (m *Middlewares) JWT(next http.Handler) http.Handler {
	m.log.Debug("Cookie request enrichment with an identifier is enabled")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := m.log.With("request_id", mwChi.GetReqID(r.Context()))

		token, err := jwt.GetFromCookie(r)
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				log.Warn(err.Error())
			}
			generateJWT(log, next, w, r)
			return
		}

		if valid := jwt.TokenValidate(token); !valid {
			generateJWT(log, next, w, r)
			return
		}

		ctx := context.WithValue(r.Context(), jwt.ContextKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateJWT(
	log *logger.Logger,
	next http.Handler,
	w http.ResponseWriter,
	r *http.Request,
) {
	token, err := jwt.GenerateToken()
	if err != nil {
		log.Debugf("failed to generate jwt token: %s", err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  jwt.CookieKey,
		Value: token,
	})

	ctx := context.WithValue(r.Context(), jwt.ContextKey, token)
	next.ServeHTTP(w, r.WithContext(ctx))
}
