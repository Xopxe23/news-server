package rest

import "net/http"

type ArticlesService interface {
}

func (h *Handler) getAllArticles(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) createArticle(w http.ResponseWriter, r *http.Request)  {}
func (h *Handler) getArticleById(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) updateArticle(w http.ResponseWriter, r *http.Request)  {}
func (h *Handler) deleteArticle(w http.ResponseWriter, r *http.Request)  {}
