package services

import (
	"context"

	r "github.com/AnatolyGolang/book-shop/internal/app/repositories"
	"github.com/AnatolyGolang/book-shop/internal/app/services/models"
)

type CategoryServiceImpl struct {
	repository r.CategoryRepository
}

func NewCategoryService(repo r.CategoryRepository) *CategoryServiceImpl {
	return &CategoryServiceImpl{
		repository: repo,
	}
}

func (s *CategoryServiceImpl) GetCategory(ctx context.Context, id int) (models.DomainCategory, error) {
	category, err := s.repository.GetCategory(ctx, id)
	if err != nil {
		return models.DomainCategory{}, err
	}
	return models.ToDomainCategory(category), nil
}

func (s *CategoryServiceImpl) CreateCategory(ctx context.Context, category models.DomainCategory) (models.DomainCategory, error) {
	categoryRep, err := s.repository.CreateCategory(ctx, category)
	if err != nil {
		return models.DomainCategory{}, err
	}
	return models.ToDomainCategory(categoryRep), nil
}

func (s *CategoryServiceImpl) UpdateCategory(ctx context.Context, id int, category models.DomainCategory) (models.DomainCategory, error) {
	bookRep, err := s.repository.UpdateCategory(ctx, id, category)
	if err != nil {
		return models.DomainCategory{}, err
	}

	return models.ToDomainCategory(bookRep), nil
}

func (s *CategoryServiceImpl) DeleteCategory(ctx context.Context, id int) error {
	err := s.repository.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *CategoryServiceImpl) GetCategories(ctx context.Context) ([]models.DomainCategory, error) {
	categories, err := s.repository.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	var domainCategories []models.DomainCategory
	for _, category := range categories {
		domainCategories = append(domainCategories, models.ToDomainCategory(category))
	}

	return domainCategories, nil
}
