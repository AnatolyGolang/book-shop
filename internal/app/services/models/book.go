package models

import (
	"book-shop/internal/app/repositories/models"
	"time"
)

type DomainBook struct {
	ID         int
	Title      string
	Year       int
	Author     string
	Price      int
	Amount     int
	CategoryID int
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

func ToDomainBook(b models.Book) DomainBook {
	return DomainBook{
		ID:         b.ID,
		Title:      b.Title,
		Year:       b.Year,
		Author:     b.Author,
		Price:      b.Price,
		Amount:     b.Amount,
		CategoryID: b.CategoryID,
		CreatedAt:  b.CreatedAt,
		UpdatedAt:  b.UpdatedAt,
	}
}
