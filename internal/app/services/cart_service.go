package services

import (
	"context"
	"fmt"
	"log"
	"time"

	r "github.com/AnatolyGolang/book-shop/internal/app/repositories"
)

type CartServiceImpl struct {
	repository r.CartRepository
}

func NewCartService(repo r.CartRepository) *CartServiceImpl {
	return &CartServiceImpl{repository: repo}
}

func (s *CartServiceImpl) UpdateCart(ctx context.Context, userID int, bookIds []int) error {
	err := s.repository.UpdateCart(ctx, userID, bookIds)
	if err != nil {
		return fmt.Errorf("error adding books to cart: %w", err)
	}
	return nil
}

func (s *CartServiceImpl) CartCleanupScheduler() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			<-ticker.C
			err := s.repository.CleanupExpiredCartItems(context.Background())
			if err != nil {
				log.Printf("cart cleanup error: %v", err)
			}
		}
	}()
}
