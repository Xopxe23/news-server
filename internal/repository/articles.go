package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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

func (r *ArticlesRepository) GetById(ctx context.Context, articleId int) (domain.ArticleOutput, error) {
	var article domain.ArticleOutput
	query := `SELECT ar.id, CONCAT(au.name, ' ', au.surname) as author, ar.title, ar.content, ar.created_at 
			  FROM articles ar INNER JOIN authors au ON ar.author_id = au.id WHERE ar.id = $1;`
	err := r.db.QueryRow(query, articleId).Scan(&article.Id, &article.Author, &article.Title, &article.Content, &article.CreatedAt)

	return article, err
}

func (r *ArticlesRepository) AddInBookmars(ctx context.Context, id, userId int) error {
	_, err := r.db.Exec("INSERT INTO bookmarks (user_id, article_id) VALUES ($1, $2)", userId, id)
	return err
}

func (r *ArticlesRepository) Update(ctx context.Context, articleId int, input domain.UpdateArticleInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	if input.Content != nil {
		setValues = append(setValues, fmt.Sprintf("content = $%d", argId))
		args = append(args, *input.Content)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE articles SET %s WHERE id = $%d", setQuery, argId)
	fmt.Println(query)
	args = append(args, articleId)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *ArticlesRepository) Delete(ctx context.Context, articleId int) error {
	_, err := r.db.Exec("DELETE FROM articles WHERE id = $1", articleId)
	return err
}
