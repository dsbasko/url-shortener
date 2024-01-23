package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/dsbasko/yandex-go-shortener/internal/jwt"
	"github.com/go-chi/chi/v5/middleware"
)

func (m *Middlewares) JWT(next http.Handler) http.Handler {
	m.log.Debug("Cookie request enrichment with an identifier is enabled")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := m.log.With("request_id", middleware.GetReqID(r.Context()))

		generate := func() {
			token, err := jwt.GenerateToken()
			if err != nil {
				log.Error(err.Error())
			}

			http.SetCookie(w, &http.Cookie{
				Name:  jwt.CookieKey,
				Value: token,
			})

			ctx := context.WithValue(r.Context(), jwt.ContextKey, token)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		token, err := jwt.GetFromCookie(r)
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				log.Warn(err.Error())
			}
			generate()
			return
		}

		valid := jwt.TokenValidate(token)
		if !valid {
			generate()
			return
		}

		ctx := context.WithValue(r.Context(), jwt.ContextKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
