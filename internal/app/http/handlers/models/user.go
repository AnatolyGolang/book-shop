package models

import "book-shop/internal/app/services/models"

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ToDomainUser(email string, password string) models.DomainUser {
	return models.DomainUser{
		Email:    email,
		Password: password,
	}
}
