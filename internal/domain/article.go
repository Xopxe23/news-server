package domain

import "time"

type Article struct {
	Id        int       `json:"id"`
	AuthorId  int       `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
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
