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
