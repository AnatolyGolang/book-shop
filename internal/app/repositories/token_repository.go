package repositories

import (
	"context"
	"fmt"
	"time"

	"book-shop/internal/pkg/postgres"
)

type TokenRepositoryImpl struct {
	db *postgres.DBConnection
}

func NewTokenRepository(db *postgres.DBConnection) *TokenRepositoryImpl {
	return &TokenRepositoryImpl{db: db}
}

func (r *TokenRepositoryImpl) SaveToken(ctx context.Context, userId int, token string, expiresAt time.Time) error {
	query := `INSERT INTO user_tokens(user_id, token, expires_at)
			  VALUES($1, $2, $3)`

	_, err := r.db.Exec(ctx, query, userId, token, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}
	return nil
}

func (r *TokenRepositoryImpl) DeleteToken(ctx context.Context, token string) error {
	query := `DELETE FROM user_tokens WHERE token = $1`
	result, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("token not found")
	}
	return nil
}

func (r *TokenRepositoryImpl) CleanupExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM user_tokens WHERE expires_at < now()`
	_, err := r.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}
	return nil
}
