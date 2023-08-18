package service

import (
	"context"
	"time"

	"github.com/xopxe23/news-server/internal/domain"
)

type ArticlesRepository interface {
	Create(сtx context.Context, input domain.Article) (int, error)
	GetAll(ctx context.Context) ([]domain.ArticleOutput, error)
}

type AuthorsRepository interface {
	Create(сtx context.Context, author domain.Author) (int, error)
	GetAll(ctx context.Context) ([]domain.Author, error)
	GetById(ctx context.Context, id int) (domain.Author, error)
	Update(ctx context.Context, id int, input domain.UpdateAuthorInput) error
	Delete(ctx context.Context, id int) error
}

type ArticlesService struct {
	authorsRepo  AuthorsRepository
	articlesRepo ArticlesRepository
}

func NewArticlesService(authorsRepo AuthorsRepository, articlesRepo ArticlesRepository) *ArticlesService {
	return &ArticlesService{
		authorsRepo:  authorsRepo,
		articlesRepo: articlesRepo,
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

func (s *ArticlesService) UpdateAuthor(ctx context.Context, id int, input domain.UpdateAuthorInput) error {
	return s.authorsRepo.Update(ctx, id, input)
}

func (s *ArticlesService) DeleteAuthor(ctx context.Context, id int) error {
	return s.authorsRepo.Delete(ctx, id)
}

func (s *ArticlesService) CreateArticle(ctx context.Context, input domain.Article) (int, error) {
	if input.CreatedAt.IsZero() {
		input.CreatedAt = time.Now()
	}

	return s.articlesRepo.Create(ctx, input)
}

func (s *ArticlesService) GetAllArticles(ctx context.Context) ([]domain.ArticleOutput, error) {
	return s.articlesRepo.GetAll(ctx)
}
