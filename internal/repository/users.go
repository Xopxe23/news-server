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
