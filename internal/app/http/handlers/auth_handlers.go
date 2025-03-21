package handlers

import (
	"book-shop/internal/app/http/handlers/errors"
	models2 "book-shop/internal/app/http/handlers/models"
	"book-shop/internal/app/utils"
	"encoding/json"
	"net/http"
	"strings"
)

func (h HttpServer) SignUp(w http.ResponseWriter, r *http.Request) {
	var authRequest models2.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		errors.BadRequest("invalid-json", err, w, r)
		return
	}

	if err := authRequest.Validate(); err != nil {
		errors.BadRequest("invalid-request", err, w, r)
		return
	}

	hashedPassword, err := utils.GetHash(authRequest.Password)
	if err != nil {
		errors.RespondWithError(err, w, r)
		return
	}

	_, err = h.userService.CreateUser(r.Context(), models2.ToDomainUser(authRequest.Email, hashedPassword))
	if err != nil {
		errors.RespondWithError(err, w, r)
		return
	}

	errors.RespondOK(map[string]bool{"ok": true}, w)
}

func (h HttpServer) SignIn(w http.ResponseWriter, r *http.Request) {
	var authRequest models2.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		errors.BadRequest("invalid-json", err, w, r)
		return
	}

	if err := authRequest.Validate(); err != nil {
		errors.BadRequest("invalid-request", err, w, r)
		return
	}

	user, err := h.userService.GetUserByName(r.Context(), authRequest.Email)
	if err != nil {
		errors.RespondWithError(err, w, r)
		return
	}

	if !utils.CheckHash(authRequest.Password, user.Password) {
		errors.BadRequest("invalid-password", nil, w, r)
		return
	}

	token, err := h.jwtService.GenerateJWT(r.Context(), user)
	if err != nil {
		errors.RespondWithError(err, w, r)
		return
	}

	errors.RespondOK(map[string]string{"token": token}, w)
}

func (h HttpServer) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get(AuthorizationHeader)
	if authHeader == "" {
		errors.BadRequest("missing-token", nil, w, r)
		return
	}

	if !strings.HasPrefix(authHeader, BearerPrefix) {
		errors.BadRequest("invalid-token-format", nil, w, r)
		return
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, BearerPrefix))
	if token == "" {
		errors.BadRequest("empty-token", nil, w, r)
		return
	}

	err := h.jwtService.RevokeToken(r.Context(), token)
	if err != nil {
		errors.RespondWithError(err, w, r)
		return
	}

	errors.RespondOK(map[string]string{"message": "Logged out successfully"}, w)
}
