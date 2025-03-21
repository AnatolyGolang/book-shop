package models

import (
	"time"

	"github.com/AnatolyGolang/book-shop/internal/app/repositories/models"
)

type DomainCategory struct {
	Id        int
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func ToDomainCategory(c models.Category) DomainCategory {
	return DomainCategory{
		Id:        c.Id,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
