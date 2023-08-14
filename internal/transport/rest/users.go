package rest

import "net/http"

type UsersService interface {
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {}