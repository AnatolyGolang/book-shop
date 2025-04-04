package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AnatolyGolang/book-shop/internal/app/http/handlers/errors"
	"github.com/AnatolyGolang/book-shop/internal/app/http/handlers/models"
	"github.com/AnatolyGolang/book-shop/internal/app/utils"
)

func (h HttpServer) SignUp(w http.ResponseWriter, r *http.Request) {
	var authRequest models.AuthRequest
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

	_, err = h.userService.CreateUser(r.Context(), models.ToDomainUser(authRequest.Email, hashedPassword))
	if err != nil {
		errors.RespondWithError(err, w, r)
		return
	}

	errors.RespondOK(map[string]bool{"ok": true}, w)
}

func (h HttpServer) SignIn(w http.ResponseWriter, r *http.Request) {
	var authRequest models.AuthRequest
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
	authHeader := r.Header.Get(utils.AuthorizationHeader)
	if authHeader == "" {
		errors.BadRequest("missing-token", nil, w, r)
		return
	}

	if !strings.HasPrefix(authHeader, utils.BearerPrefix) {
		errors.BadRequest("invalid-token-format", nil, w, r)
		return
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, utils.BearerPrefix))
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
