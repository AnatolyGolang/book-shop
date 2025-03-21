package handlers

import (
	errors2 "book-shop/internal/app/http/handlers/errors"
	"book-shop/internal/app/http/handlers/models"
	se "book-shop/internal/app/services/errors"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h HttpServer) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req models.CategoryCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors2.BadRequest("invalid-request-body", err, w, r)
		return
	}

	if err := req.Validate(); err != nil {
		errors2.BadRequest("validation-error", err, w, r)
		return
	}

	category, err := h.categoryService.CreateCategory(r.Context(), models.ToServiceCategoryCreate(req))
	if err != nil {
		if errors.Is(err, se.ErrRequired) {
			errors2.NotFound("all field need to be filled", err, w, r)
			return
		}
		errors2.RespondWithError(err, w, r)
		return
	}

	response := models.ToCategoryResponse(category)
	errors2.RespondCreated(response, w)
}

func (h HttpServer) GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryId, err := strconv.Atoi(vars["category_id"])
	if err != nil {
		errors2.BadRequest("invalid-category-id", err, w, r)
		return
	}
	category, err := h.categoryService.GetCategory(r.Context(), categoryId)
	if err != nil {
		if errors.Is(err, se.ErrNotFound) {
			errors2.NotFound("category-not-found", err, w, r)
			return
		}
		errors2.RespondWithError(err, w, r)
		return
	}

	response := models.ToCategoryResponse(category)

	errors2.RespondOK(response, w)
}

func (h HttpServer) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryId, err := strconv.Atoi(vars["category_id"])
	if err != nil {
		errors2.BadRequest("invalid-category-id", err, w, r)
		return
	}

	var req models.CategoryUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors2.BadRequest("invalid-request-body", err, w, r)
		return
	}

	if err := req.Validate(); err != nil {
		errors2.BadRequest("validation-error", err, w, r)
		return
	}

	category, err := h.categoryService.UpdateCategory(r.Context(), categoryId, models.ToServiceCategoryUpdate(req))
	if err != nil {
		if errors.Is(err, se.ErrNotFound) {
			errors2.NotFound("category-not-found", err, w, r)
			return
		}
		errors2.RespondWithError(err, w, r)
		return
	}

	response := models.ToCategoryResponse(category)
	errors2.RespondOK(response, w)
}

func (h HttpServer) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryId, err := strconv.Atoi(vars["category_id"])
	if err != nil {
		errors2.BadRequest("invalid-category-id", err, w, r)
		return
	}

	err = h.categoryService.DeleteCategory(r.Context(), categoryId)
	if err != nil {
		if errors.Is(err, se.ErrNotFound) {
			errors2.NotFound("book-not-found", err, w, r)
			return
		}
		errors2.RespondWithError(err, w, r)
		return
	}

	errors2.RespondNoContent(w)
}

func (h HttpServer) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryService.GetCategories(r.Context())
	if err != nil {
		if errors.Is(err, se.ErrNotFound) {
			errors2.NotFound("categories-not-found", err, w, r)
			return
		}
		errors2.RespondWithError(err, w, r)
		return
	}

	response := models.ToCategoriesResponse(categories)

	errors2.RespondOK(response, w)
}
