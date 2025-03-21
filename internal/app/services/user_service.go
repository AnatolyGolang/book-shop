package services

import (
	"context"

	r "github.com/AnatolyGolang/book-shop/internal/app/repositories"
	"github.com/AnatolyGolang/book-shop/internal/app/services/models"
)

type UserServiceImpl struct {
	repository r.UserRepository
}

func NewUserService(repo r.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		repository: repo,
	}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, domainUser models.DomainUser) (models.DomainUser, error) {
	book, err := s.repository.CreateUser(ctx, domainUser)
	if err != nil {
		return models.DomainUser{}, err
	}
	return models.ToDomainUser(book), nil
}

func (s *UserServiceImpl) GetUserByName(ctx context.Context, name string) (models.DomainUser, error) {
	user, err := s.repository.GetUserByEmail(ctx, name)
	if err != nil {
		return models.DomainUser{}, err
	}
	return models.ToDomainUser(user), nil
}

func (s *UserServiceImpl) GetUserById(ctx context.Context, id int) (models.DomainUser, error) {
	user, err := s.repository.GetUserById(ctx, id)
	if err != nil {
		return models.DomainUser{}, err
	}
	return models.ToDomainUser(user), nil
}
