package handlers

import (
	"encoding/json"
	"net/http"

	he "github.com/AnatolyGolang/book-shop/internal/app/http/handlers/errors"
	"github.com/AnatolyGolang/book-shop/internal/app/http/handlers/models"
	se "github.com/AnatolyGolang/book-shop/internal/app/services/models"
)

func (h HttpServer) AddToCart(w http.ResponseWriter, r *http.Request) {
	user, err := se.GetUserFromContext(r.Context())
	if err != nil {
		he.Unauthorised("unauthorized", err, w, r)
		return
	}

	var cartReq models.CartUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&cartReq); err != nil {
		he.BadRequest("invalid-request-body", err, w, r)
		return
	}

	if len(cartReq.BookIds) <= 0 {
		he.BadRequest("invalid-request-body", err, w, r)
		return
	}

	err = h.cartService.UpdateCart(r.Context(), user.Id, cartReq.BookIds)
	if err != nil {
		he.RespondWithError(err, w, r)
		return
	}

	he.RespondOK(map[string]string{"message": "Book added to cart"}, w)
}
