package services

import (
	"context"

	"github.com/AnatolyGolang/book-shop/internal/app/services/models"
)

type BookService interface {
	GetBook(ctx context.Context, id int) (models.DomainBook, error)
	CreateBook(ctx context.Context, book models.DomainBook) (models.DomainBook, error)
	UpdateBook(ctx context.Context, id int, book models.DomainBook) (models.DomainBook, error)
	DeleteBook(ctx context.Context, id int) error
	GetBooksByCategories(ctx context.Context, categoryIDs []int, limit int, offset int) ([]models.DomainBook, int, error)
}

type CategoryService interface {
	GetCategory(ctx context.Context, id int) (models.DomainCategory, error)
	CreateCategory(ctx context.Context, category models.DomainCategory) (models.DomainCategory, error)
	UpdateCategory(ctx context.Context, id int, category models.DomainCategory) (models.DomainCategory, error)
	DeleteCategory(ctx context.Context, id int) error
	GetCategories(ctx context.Context) ([]models.DomainCategory, error)
}

type UserService interface {
	CreateUser(ctx context.Context, user models.DomainUser) (models.DomainUser, error)
	GetUserByName(ctx context.Context, name string) (models.DomainUser, error)
	GetUserById(ctx context.Context, id int) (models.DomainUser, error)
}

type JWTService interface {
	GenerateJWT(ctx context.Context, user models.DomainUser) (string, error)
	GetUser(ctx context.Context, token string) (models.DomainUser, error)
	RevokeToken(ctx context.Context, token string) error
	StartTokenCleanupScheduler()
}

type CartService interface {
	UpdateCart(ctx context.Context, userID int, bookIds []int) error
	CartCleanupScheduler()
}
