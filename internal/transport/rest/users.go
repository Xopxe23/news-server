package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/xopxe23/news-server/internal/domain"
)

type UsersService interface {
	SignUp(ctx context.Context, input domain.SignUpInput) error
	SignIn(ctx context.Context, input domain.SignInInput) (string, string, error)
	RefreshTokens(ctx context.Context, token string) (string, string, error)
}

// @Summary Sign Up
// @Tags Users auth
// @ID sign-up
// @Accept json
// @Produce json
// @Param input body domain.SignUpInput true "Sign up input"
// @Success 200
// @Failure 400
// @Router /auth/sign-up [post]
func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var input domain.SignUpInput
	if err := json.Unmarshal(reqBytes, &input); err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := input.Validate(); err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.usersService.SignUp(r.Context(), input); err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(map[string]string{
		"status": "success",
	})
	if err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Sign In
// @Tags Users auth
// @ID sign-in
// @Accept json
// @Produce json
// @Param input body domain.SignInInput true "Sign in input"
// @Success 200 {string} string
// @Failure 400
// @Failure 500
// @Router /auth/sign-in [post]
func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var input domain.SignInInput
	if err := json.Unmarshal(reqBytes, &input); err != nil {
		logError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := input.Validate(); err != nil {
		logError("sigIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.usersService.SignIn(r.Context(), input)
	if err != nil {
		logError("signIn", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]string{
		"token": accessToken,
	})
	if err != nil {
		logError("signIn", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Refresh
// @Tags Users auth
// @ID refresh
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Failure 500
// @Router /auth/refresh [get]
func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		logError("refresh", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logrus.Infof("%s", cookie.Value)

	accessToken, refreshToken, err := h.usersService.RefreshTokens(r.Context(), cookie.Value)
	if err != nil {
		logError("refresh", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]string{
		"token": accessToken,
	})
	if err != nil {
		logError("refresh", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}
