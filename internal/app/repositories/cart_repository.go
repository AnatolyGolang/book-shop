package repositories

import "book-shop/internal/pkg/postgres"

type CartRepositoryImpl struct {
	db *postgres.DBConnection
}

func NewCartRepository(db *postgres.DBConnection) *CartRepositoryImpl {
	return &CartRepositoryImpl{db: db}
}
