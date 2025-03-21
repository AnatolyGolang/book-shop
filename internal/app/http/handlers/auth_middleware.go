package handlers

import (
	he "book-shop/internal/app/http/handlers/errors"

	"context"
	"net/http"
	"strings"
)

func (h HttpServer) CheckAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(AuthorizationHeader)
		token = strings.TrimSpace(strings.TrimPrefix(token, BearerPrefix))
		user, err := h.jwtService.GetUser(r.Context(), token)

		if err != nil {
			he.InternalError("invalid-token", err, w, r)
			return
		}

		if user.Email == "" {
			he.InternalError("empty email in token", nil, w, r)
			return
		}

		if !user.IsAdmin {
			he.Unauthorised("user not admin", nil, w, r)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserKey, user)
		next(w, r.WithContext(ctx))
	}
}

func (h HttpServer) CheckAuthorizedUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(AuthorizationHeader)
		token = strings.TrimSpace(strings.TrimPrefix(token, BearerPrefix))
		user, err := h.jwtService.GetUser(r.Context(), token)

		if err != nil {
			he.InternalError("invalid-token", err, w, r)
			return
		}

		if user.Email == "" {
			he.InternalError("empty email in token", nil, w, r)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserKey, user)
		next(w, r.WithContext(ctx))
	}
}
