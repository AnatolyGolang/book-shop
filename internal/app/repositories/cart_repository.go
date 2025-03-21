package repositories

import (
	"book-shop/internal/app/logger"
	"book-shop/internal/pkg/postgres"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type CartRepositoryImpl struct {
	db *postgres.DBConnection
}

func NewCartRepository(db *postgres.DBConnection) *CartRepositoryImpl {
	return &CartRepositoryImpl{db: db}
}

func (r *CartRepositoryImpl) UpdateCart(ctx context.Context, userID int, bookIds []int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			logger.Logger.Error("failed to rollback", zap.Error(err))
		}
	}(tx, ctx)

	query := `SELECT id, amount FROM books WHERE id = ANY($1) FOR UPDATE`
	rows, err := tx.Query(ctx, query, bookIds)
	if err != nil {
		return fmt.Errorf("failed to lock amounts: %w", err)
	}
	defer rows.Close()

	amountMap := make(map[int]int)
	for rows.Next() {
		var bookID, amount int
		if err := rows.Scan(&bookID, &amount); err != nil {
			return fmt.Errorf("failed to scan amount: %w", err)
		}
		amountMap[bookID] = amount
	}

	for _, bookID := range bookIds {
		if amountMap[bookID] == 0 {
			return fmt.Errorf("book out of stock %d", bookID)
		}
	}

	query = `UPDATE books SET amount = amount - 1, updated_at = NOW() WHERE id = ANY($1)`
	_, err = tx.Exec(ctx, query, bookIds)
	if err != nil {
		return fmt.Errorf("failed to reduce amount: %w", err)
	}

	query = `
        INSERT INTO carts (user_id, book_ids, updated_at)
        VALUES ($1, $2, NOW())
        ON CONFLICT (user_id)
        DO UPDATE SET book_ids = ARRAY(
            SELECT DISTINCT unnest(carts.book_ids || EXCLUDED.book_ids)
        ), updated_at = NOW()`
	_, err = tx.Exec(ctx, query, userID, bookIds)
	if err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *CartRepositoryImpl) CleanupExpiredCartItems(ctx context.Context) error {
	query := `
        UPDATE carts
        SET book_ids = ARRAY(
            SELECT unnest(book_ids)
            FROM carts
            WHERE updated_at > NOW() - INTERVAL '30 minutes'
        ), updated_at = NOW()
        WHERE updated_at <= NOW() - INTERVAL '30 minutes'`

	_, err := r.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to remove expired cart items: %w", err)
	}
	return nil
}
