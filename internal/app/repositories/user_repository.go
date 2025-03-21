package repositories

import (
	rm "book-shop/internal/app/repositories/models"
	sm "book-shop/internal/app/services/models"
	"book-shop/internal/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type UserRepositoryImpl struct {
	db *postgres.DBConnection
}

func NewUserRepository(db *postgres.DBConnection) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (rm.User, error) {
	if email == "" {
		return rm.User{}, fmt.Errorf("email can not be empty")
	}

	var user rm.User
	query := `SELECT id, email, password, is_admin, created_at, updated_at 
              	FROM users
              WHERE email = $1`

	row := r.db.QueryRow(ctx, query, email)
	err := row.Scan(
		&user.Id, &user.Email, &user.Password, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return rm.User{}, fmt.Errorf("user not found")
		}
		return rm.User{}, fmt.Errorf("failed to get a user: %w", err)
	}

	return user, nil
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user sm.DomainUser) (rm.User, error) {
	query := `INSERT INTO users (email, password) 
              VALUES ($1, $2) 
              RETURNING id, email, password, is_admin, created_at, updated_at`

	var newUser rm.User
	err := r.db.QueryRow(ctx, query, user.Email, user.Password).Scan(
		&newUser.Id, &newUser.Email, &newUser.Password, &newUser.IsAdmin,
		&newUser.CreatedAt, &newUser.UpdatedAt,
	)
	if err != nil {
		return rm.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

func (r *UserRepositoryImpl) GetUserById(ctx context.Context, id int) (rm.User, error) {
	if id == 0 {
		return rm.User{}, fmt.Errorf("id can not be 0")
	}

	var user rm.User
	query := `SELECT id, email, password, is_admin, created_at, updated_at 
              	FROM users
              WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(
		&user.Id, &user.Email, &user.Password, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return rm.User{}, fmt.Errorf("user not found")
		}
		return rm.User{}, fmt.Errorf("failed to get a user: %w", err)
	}

	return user, nil
}
