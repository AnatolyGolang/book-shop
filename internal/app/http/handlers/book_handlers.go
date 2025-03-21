package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	he "github.com/AnatolyGolang/book-shop/internal/app/http/handlers/errors"
	"github.com/AnatolyGolang/book-shop/internal/app/http/handlers/models"
	se "github.com/AnatolyGolang/book-shop/internal/app/services/errors"

	"github.com/gorilla/mux"
)

const (
	MinLimit = 50
	MaxLimit = 100
)

func (h HttpServer) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["book_id"])
	if err != nil {
		he.BadRequest("invalid-book-id", err, w, r)
		return
	}
	book, err := h.bookService.GetBook(r.Context(), bookID)
	if err != nil {
		if errors.Is(err, se.ErrNotFound) {
			he.NotFound("book-not-found", err, w, r)
			return
		}
		he.RespondWithError(err, w, r)
		return
	}

	response := models.ToBookResponse(book)

	he.RespondOK(response, w)
}

func (h HttpServer) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req models.BookCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		he.BadRequest("invalid-request-body", err, w, r)
		return
	}

	if err := req.Validate(); err != nil {
		he.BadRequest("validation-error", err, w, r)
		return
	}

	book, err := h.bookService.CreateBook(r.Context(), models.ToServiceBookCreate(req))
	if err != nil {
		if errors.Is(err, se.ErrRequired) {
			he.NotFound("all field need to be filled", err, w, r)
			return
		}
		he.RespondWithError(err, w, r)
		return
	}

	response := models.ToBookResponse(book)
	he.RespondCreated(response, w)
}

func (h HttpServer) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["book_id"])
	if err != nil {
		he.BadRequest("invalid-book-id", err, w, r)
		return
	}

	var req models.BookUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		he.BadRequest("invalid-request-body", err, w, r)
		return
	}

	if err := req.Validate(); err != nil {
		he.BadRequest("validation-error", err, w, r)
		return
	}

	book, err := h.bookService.UpdateBook(r.Context(), bookID, models.ToServiceBookUpdate(req))
	if err != nil {
		if errors.Is(err, se.ErrNotFound) {
			he.NotFound("book-not-found", err, w, r)
			return
		}
		he.RespondWithError(err, w, r)
		return
	}

	response := models.ToBookResponse(book)
	he.RespondOK(response, w)
}

func (h HttpServer) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["book_id"])
	if err != nil {
		he.BadRequest("invalid-book-id", err, w, r)
		return
	}

	err = h.bookService.DeleteBook(r.Context(), bookID)
	if err != nil {
		if errors.Is(err, se.ErrNotFound) {
			he.NotFound("book-not-found", err, w, r)
			return
		}
		he.RespondWithError(err, w, r)
		return
	}

	he.RespondNoContent(w)
}

func (h HttpServer) GetBooks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["book_id"])
	if err != nil {
		he.BadRequest("invalid-book-id", err, w, r)
		return
	}

	err = h.bookService.DeleteBook(r.Context(), bookID)
	if err != nil {
		if errors.Is(err, se.ErrNotFound) {
			he.NotFound("book-not-found", err, w, r)
			return
		}
		he.RespondWithError(err, w, r)
		return
	}

	he.RespondNoContent(w)
}

func (h HttpServer) GetBooksByCategories(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	categoryIDsReq, ok := queryParams["category_ids"]
	if !ok || len(categoryIDsReq) == 0 {
		he.BadRequest("missing-category-ids", errors.New("at least one categoryId must be provided"), w, r)
		return
	}

	var categoryIDs []int
	for _, idStr := range categoryIDsReq {
		ids := strings.Split(idStr, ",")
		for _, id := range ids {
			categoryID, err := strconv.Atoi(strings.TrimSpace(id))
			if err != nil {
				he.BadRequest("invalid-category-id", err, w, r)
				return
			}
			categoryIDs = append(categoryIDs, categoryID)
		}
	}

	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit < MinLimit {
		limit = MinLimit
	} else if limit > MaxLimit {
		limit = MaxLimit
	}

	offset := (page - 1) * limit

	books, total, err := h.bookService.GetBooksByCategories(r.Context(), categoryIDs, limit, offset)
	if err != nil {
		if errors.Is(err, se.ErrNotFound) {
			he.NotFound("books-not-found", err, w, r)
			return
		}
		he.RespondWithError(err, w, r)
		return
	}

	response := models.PaginationResponse{
		Books: models.ToBooksResponse(books),
		Meta: models.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}

	he.RespondOK(response, w)
}
