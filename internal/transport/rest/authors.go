package rest

import "net/http"

type AuthorsService interface{}

func (h *Handler) getAllAuthors(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) createAuthor(w http.ResponseWriter, r *http.Request)  {}
func (h *Handler) getAuthorById(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) updateAuthor(w http.ResponseWriter, r *http.Request)  {}
func (h *Handler) deleteAuthor(w http.ResponseWriter, r *http.Request)  {}
