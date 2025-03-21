package postgres

import (
	"book-shop/internal/app/logger"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConnection struct {
	*pgxpool.Pool
}

func Dial(ctx context.Context, dsn string) (*DBConnection, error) {
	dbPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("error connection to DB: %w", err)
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("can not PING db: %w", err)
	}

	logger.Logger.Info("successfully connected to DB")
	return &DBConnection{dbPool}, nil
}
