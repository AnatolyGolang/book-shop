package models

import (
	"context"
	"fmt"
	"time"

	"github.com/AnatolyGolang/book-shop/internal/app/repositories/models"
	"github.com/AnatolyGolang/book-shop/internal/app/utils"
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

func GetUserFromContext(ctx context.Context) (DomainUser, error) {
	contextUser := ctx.Value(utils.ContextUserKey)
	if contextUser == nil {
		return DomainUser{}, fmt.Errorf("no user in context")
	}
	user, ok := contextUser.(DomainUser)
	if !ok {
		return DomainUser{}, fmt.Errorf("no user in context")
	}
	return user, nil
}
