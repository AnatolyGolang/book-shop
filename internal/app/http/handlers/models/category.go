package models

import (
	"book-shop/internal/app/services/models"
	"fmt"
)

type CategoryCreateRequest struct {
	Name string `json:"name"`
}

func (cr *CategoryCreateRequest) Validate() error {
	if cr.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

func (cr *CategoryUpdateRequest) Validate() error {
	if cr.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

func ToServiceCategoryCreate(request CategoryCreateRequest) models.DomainCategory {
	return models.DomainCategory{
		Name: request.Name,
	}
}

func ToServiceCategoryUpdate(request CategoryUpdateRequest) models.DomainCategory {
	return models.DomainCategory{
		Name: request.Name,
	}
}

func ToCategoriesResponse(categories []models.DomainCategory) []CategoryResponse {
	response := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		response[i] = ToCategoryResponse(category)
	}
	return response
}

func ToCategoryResponse(b models.DomainCategory) CategoryResponse {
	return CategoryResponse{
		Id:   b.Id,
		Name: b.Name,
	}
}

type CategoryResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CategoryUpdateRequest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
