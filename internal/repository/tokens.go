package repository

import (
	"context"
	"database/sql"

	"github.com/xopxe23/news-server/internal/domain"
)

type TokensRepository struct {
	db *sql.DB
}

func NewTokensRepository(db *sql.DB) *TokensRepository {
	return &TokensRepository{db: db}
}

func (r *TokensRepository) Create(ctx context.Context, token domain.RefreshSession) error {
	_, err := r.db.Exec("INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		token.UserId, token.Token, token.ExpiresAt)
	return err
}

func (r *TokensRepository) GetToken(ctx context.Context, token string) (domain.RefreshSession, error) {
	var t domain.RefreshSession
	err := r.db.QueryRow("SELECT * FROM refresh_tokens WHERE token = $1", token).Scan(&t.Id, &t.UserId, &t.Token, &t.ExpiresAt)
	if err != nil {
		return t, err
	}
	_, err = r.db.Exec("DELETE FROM refresh_tokens WHERE user_id = $1", t.UserId)
	return t, err
}
