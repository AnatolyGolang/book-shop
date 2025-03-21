package repositories

import (
	"context"
	"time"

	"github.com/AnatolyGolang/book-shop/internal/app/repositories/models"
	domain "github.com/AnatolyGolang/book-shop/internal/app/services/models"
)

type BookRepository interface {
	GetBook(ctx context.Context, id int) (models.Book, error)
	CreateBook(ctx context.Context, book domain.DomainBook) (models.Book, error)
	UpdateBook(ctx context.Context, id int, book domain.DomainBook) (models.Book, error)
	DeleteBook(ctx context.Context, id int) error
	GetBooksByCategories(ctx context.Context, categoryIDs []int, limit int, offset int) ([]models.Book, int, error)
}

type CategoryRepository interface {
	GetCategory(ctx context.Context, id int) (models.Category, error)
	CreateCategory(ctx context.Context, category domain.DomainCategory) (models.Category, error)
	UpdateCategory(ctx context.Context, id int, category domain.DomainCategory) (models.Category, error)
	DeleteCategory(ctx context.Context, id int) error
	GetCategories(ctx context.Context) ([]models.Category, error)
}

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, user domain.DomainUser) (models.User, error)
	GetUserById(ctx context.Context, id int) (models.User, error)
}

type CartRepository interface {
	CleanupExpiredCartItems(ctx context.Context) error
	UpdateCart(ctx context.Context, userID int, bookIds []int) error
}

type TokenRepository interface {
	SaveToken(ctx context.Context, userId int, token string, expiresAt time.Time) error
	DeleteToken(ctx context.Context, token string) error
	CleanupExpiredTokens(ctx context.Context) error
}
