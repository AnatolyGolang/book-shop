package services

import (
	r "book-shop/internal/app/repositories"
	"book-shop/internal/app/services/models"
	"context"
)

type BookServiceImpl struct {
	repository r.BookRepository
}

func NewBookService(repo r.BookRepository) *BookServiceImpl {
	return &BookServiceImpl{
		repository: repo,
	}
}

func (s *BookServiceImpl) GetBook(ctx context.Context, id int) (models.DomainBook, error) {
	book, err := s.repository.GetBook(ctx, id)
	if err != nil {
		return models.DomainBook{}, err
	}
	return models.ToDomainBook(book), nil
}

func (s *BookServiceImpl) CreateBook(ctx context.Context, domainBook models.DomainBook) (models.DomainBook, error) {
	book, err := s.repository.CreateBook(ctx, domainBook)
	if err != nil {
		return models.DomainBook{}, err
	}
	return models.ToDomainBook(book), nil
}

func (s *BookServiceImpl) UpdateBook(ctx context.Context, id int, domainBook models.DomainBook) (models.DomainBook, error) {
	book, err := s.repository.UpdateBook(ctx, id, domainBook)
	if err != nil {
		return models.DomainBook{}, err
	}

	return models.ToDomainBook(book), nil
}

func (s *BookServiceImpl) DeleteBook(ctx context.Context, id int) error {
	err := s.repository.DeleteBook(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *BookServiceImpl) GetBooksByCategories(ctx context.Context, categoryIDs []int, limit int, offset int) ([]models.DomainBook, int, error) {
	books, total, err := s.repository.GetBooksByCategories(ctx, categoryIDs, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var domainBooks []models.DomainBook
	for _, book := range books {
		domainBooks = append(domainBooks, models.ToDomainBook(book))
	}

	return domainBooks, total, nil
}
