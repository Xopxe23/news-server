package repository

import (
	"context"
	"database/sql"

	"github.com/xopxe23/news-server/internal/domain"
)

type ArticlesRepository struct {
	db *sql.DB
}

func NewArticlesRepository(db *sql.DB) *ArticlesRepository {
	return &ArticlesRepository{db: db}
}

func (r *ArticlesRepository) Create(ctx context.Context, input domain.Article) (int, error) {
	var articleId int
	err := r.db.QueryRow("INSERT INTO articles(author_id, title, content, created_at) values($1, $2, $3, $4) RETURNING id",
		input.AuthorId, input.Title, input.Content, input.CreatedAt).Scan(&articleId)

	return articleId, err
}

func (r *ArticlesRepository) GetAll(ctx context.Context) ([]domain.ArticleOutput, error) {
	var articles []domain.ArticleOutput
	query := `SELECT ar.id, CONCAT(au.name, ' ', au.surname) as author, ar.title, ar.content, ar.created_at 
			  FROM articles ar INNER JOIN authors au ON ar.author_id = au.id;`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var article domain.ArticleOutput
		if err := rows.Scan(&article.Id, &article.Author, &article.Title, &article.Content, &article.CreatedAt); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}
