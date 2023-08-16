package service

import (
	"context"

	"github.com/xopxe23/news-server/internal/domain"
)

type ArticlesRepository interface{}

type AuthorsRepository interface {
	Create(сtx context.Context, author domain.Author) (int, error)
	GetAll(ctx context.Context) ([]domain.Author, error)
	GetById(ctx context.Context, id int) (domain.Author, error)
}

type ArticlesService struct {
	authorsRepo AuthorsRepository
}

func NewArticlesService(authorsRepo AuthorsRepository) *ArticlesService {
	return &ArticlesService{
		authorsRepo: authorsRepo,
	}
}

func (s *ArticlesService) CreateAuthor(ctx context.Context, author domain.Author) (int, error) {
	return s.authorsRepo.Create(ctx, author)
}

func (s *ArticlesService) GetAllAuthors(ctx context.Context) ([]domain.Author, error) {
	return s.authorsRepo.GetAll(ctx)
}

func (s *ArticlesService) GetAuthorById(ctx context.Context, id int) (domain.Author, error) {
	return s.authorsRepo.GetById(ctx, id)
}
