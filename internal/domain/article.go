package domain

import "time"

type Article struct {
	Id        int       `json:"id"`
	AuthorId  int       `json:"author_id"`
	Title     string    `json:"title" validate:"required,gte=10"`
	Content   string    `json:"content" validate:"required,gte=20"`
	CreatedAt time.Time `json:"created_at"`
}

type Author struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type UpdateAuthorInput struct {
	Name    *string `json:"name"`
	Surname *string `json:"surname"`
}

type UpdateArticleInput struct {
	Title     *string    `json:"title" validate:"required,gte=10"`
	Content   *string    `json:"content" validate:"required,gte=20"`
}

type ArticleOutput struct {
	Id        int       `json:"id"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *Article) Validate() error {
	return validate.Struct(a)
}
