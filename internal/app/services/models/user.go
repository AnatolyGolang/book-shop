package models

import (
	"book-shop/internal/app/repositories/models"
	"time"
)

type DomainUser struct {
	Id        int
	Email     string
	Password  string
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func ToDomainUser(u models.User) DomainUser {
	return DomainUser{
		Id:        u.Id,
		Email:     u.Email,
		Password:  u.Password,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
