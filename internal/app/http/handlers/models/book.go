package models

import (
	"fmt"

	"github.com/AnatolyGolang/book-shop/internal/app/services/models"
)

type BookCreateRequest struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	Year       int    `json:"year"`
	Price      int    `json:"price"`
	CategoryId int    `json:"category_id"`
	Amount     int    `json:"amount"`
}

type BookUpdateRequest struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	Year       int    `json:"year"`
	Price      int    `json:"price"`
	CategoryId int    `json:"category_id"`
	Amount     int    `json:"amount"`
}

func (br *BookCreateRequest) Validate() error {
	if br.Title == "" {
		return fmt.Errorf("title is required")
	}
	if br.Author == "" {
		return fmt.Errorf("author is required")
	}
	if br.Year == 0 {
		return fmt.Errorf("year is required")
	}
	if br.Price == 0 {
		return fmt.Errorf("price is required")
	}
	if br.Amount == 0 {
		return fmt.Errorf("amount is required")
	}
	if br.CategoryId == 0 {
		return fmt.Errorf("category_id is required")
	}
	return nil
}

func (br *BookUpdateRequest) Validate() error {
	if br.Title == "" {
		return fmt.Errorf("title is required")
	}
	if br.Author == "" {
		return fmt.Errorf("author is required")
	}
	if br.Year == 0 {
		return fmt.Errorf("year is required")
	}
	if br.Price == 0 {
		return fmt.Errorf("price is required")
	}
	if br.CategoryId == 0 {
		return fmt.Errorf("category_id is required")
	}
	return nil
}

type BookResponse struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	Year       int    `json:"year"`
	Price      int    `json:"price"`
	CategoryId int    `json:"category_id"`
	Amount     int    `json:"amount"`
	ID         int    `json:"id"`
}

func ToServiceBookCreate(request BookCreateRequest) models.DomainBook {
	return models.DomainBook{
		Title:      request.Title,
		Year:       request.Year,
		Author:     request.Author,
		Price:      request.Price,
		CategoryID: request.CategoryId,
		Amount:     request.Amount,
	}
}

func ToServiceBookUpdate(request BookUpdateRequest) models.DomainBook {
	return models.DomainBook{
		Title:      request.Title,
		Year:       request.Year,
		Author:     request.Author,
		Price:      request.Price,
		CategoryID: request.CategoryId,
		Amount:     request.Amount,
	}
}

func ToBookResponse(b models.DomainBook) BookResponse {
	return BookResponse{
		Title:      b.Title,
		Year:       b.Year,
		Author:     b.Author,
		Price:      b.Price,
		Amount:     b.Amount,
		CategoryId: b.CategoryID,
		ID:         b.ID,
	}
}

type PaginationResponse struct {
	Books []BookResponse `json:"books"`
	Meta  PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

func ToBooksResponse(books []models.DomainBook) []BookResponse {
	response := make([]BookResponse, len(books))
	for i, book := range books {
		response[i] = ToBookResponse(book)
	}
	return response
}
