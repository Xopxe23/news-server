package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	_ "github.com/xopxe23/news-server/docs"
)

type Handler struct {
	// articlesService ArticlesService
	// authorsService AuthorsService
	usersService UsersService
}

func NewHandler(users UsersService) *Handler {
	return &Handler{
		// articlesService: articles,
		usersService: users,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.PathPrefix("/swagger").HandlerFunc(httpSwagger.WrapHandler)

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
		auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodPost)
	}

	authors := r.PathPrefix("/authors").Subrouter()
	{
		authors.HandleFunc("", h.getAllAuthors).Methods(http.MethodGet)
		authors.HandleFunc("", h.createAuthor).Methods(http.MethodPost)
		authors.HandleFunc("", h.getAuthorById).Methods(http.MethodGet)
		authors.HandleFunc("", h.updateAuthor).Methods(http.MethodPut)
		authors.HandleFunc("", h.deleteAuthor).Methods(http.MethodDelete)
	}

	articles := r.PathPrefix("/articles").Subrouter()
	{
		articles.HandleFunc("", h.getAllArticles).Methods(http.MethodGet)
		articles.HandleFunc("", h.createArticle).Methods(http.MethodPost)
		articles.HandleFunc("/{id:[0-9]+}", h.getArticleById).Methods(http.MethodGet)
		articles.HandleFunc("/{id:[0-9]+}", h.updateArticle).Methods(http.MethodPut)
		articles.HandleFunc("/{id:[0-9]+}", h.deleteArticle).Methods(http.MethodDelete)
	}
	return r
}
