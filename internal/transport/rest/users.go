package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/xopxe23/news-server/internal/domain"
)

type UsersService interface {
	SignUp(ctx context.Context, input domain.SignUpInput) error
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	var input domain.SignUpInput
	if err := json.Unmarshal(reqBytes, &input); err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := h.usersService.SignUp(context.TODO(), input); err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {}
