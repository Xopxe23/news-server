package repository

import (
	"context"
	"database/sql"

	"github.com/xopxe23/news-server/internal/domain"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) Create(ctx context.Context, user domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (name, email, password_hash) values ($1, $2, $3)",
		user.Name, user.Email, user.Password)
	return err
}

func (r *UsersRepository) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, name, email FROM users WHERE email = $1 and password_hash = $2",
		email, password).Scan(&user.Id, &user.Name, &user.Email)
	return user, err
}

func (r *UsersRepository) GetBookmarks(ctx context.Context, userId int) ([]domain.ArticleOutput, error) {
	var articles []domain.ArticleOutput
	query := `SELECT ar.id, CONCAT(au.name, ' ', au.surname) as author, ar.title, ar.content, ar.created_at 
			  FROM articles ar INNER JOIN authors au ON ar.author_id = au.id 
			  INNER JOIN bookmarks bm ON ar.id = bm.article_id WHERE bm.user_id = $1;`
	rows, err := r.db.Query(query, userId)
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
