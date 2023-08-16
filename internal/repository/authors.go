package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/xopxe23/news-server/internal/domain"
)

type AuthorsRepository struct {
	db *sql.DB
}

func NewAuthorsRepository(db *sql.DB) *AuthorsRepository {
	return &AuthorsRepository{db: db}
}

func (r *AuthorsRepository) Create(ctx context.Context, author domain.Author) (int, error) {
	var id int
	err := r.db.QueryRow("INSERT INTO authors(name, surname) VALUES($1, $2) RETURNING id",
		author.Name, author.Surname).Scan(&id)
	return id, err
}

func (r *AuthorsRepository) GetAll(ctx context.Context) ([]domain.Author, error) {
	var authors []domain.Author
	rows, err := r.db.Query("SELECT id, name, surname FROM AUTHORS")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var author domain.Author
		if err := rows.Scan(&author.Id, &author.Name, &author.Surname); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, nil
}

func (r *AuthorsRepository) GetById(ctx context.Context, id int) (domain.Author, error) {
	var author domain.Author
	err := r.db.QueryRow("SELECT id, name, surname FROM AUTHORS WHERE id = $1",
		id).Scan(&author.Id, &author.Name, &author.Surname)

	return author, err
}

func (r *AuthorsRepository) Update(ctx context.Context, id int, input domain.UpdateAuthorInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Surname != nil {
		setValues = append(setValues, fmt.Sprintf("surname=$%d", argId))
		args = append(args, *input.Surname)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE authors SET %s WHERE id = $%d", setQuery, argId)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *AuthorsRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec("DELETE FROM authors WHERE id = $1", id)
	return err
}
