package service

import (
	"context"

	"github.com/xopxe23/news-server/internal/domain"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UsersRepository interface {
	Create(ctx context.Context, user domain.User) error
}

type UsersService struct {
	repo   UsersRepository
	hasher PasswordHasher
}

func NewUsersService(repo UsersRepository, hasher PasswordHasher) *UsersService {
	return &UsersService{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *UsersService) SignUp(ctx context.Context, input domain.SignUpInput) error {
	password, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}
	user := domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: password,
	}
	return s.repo.Create(ctx, user)
}
