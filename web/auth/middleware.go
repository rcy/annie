package auth

import (
	"context"
	"goirc/db/model"
	"net/http"
	"time"
)

type keyType int

const (
	SessionKey keyType = iota
	NickKey
)

const FromCookieKey = "annie.from"

type service struct {
	Queries *model.Queries
}

func NewService(q *model.Queries) *service {
	return &service{Queries: q}
}

func (s *service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		key := ctx.Value(SessionKey).(string)
		session, err := s.Queries.NickBySession(ctx, key)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     FromCookieKey,
				Value:    r.URL.Path,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				Expires:  time.Now().Add(time.Hour),
			})
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		ctx = context.WithValue(ctx, NickKey, session.Nick)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
