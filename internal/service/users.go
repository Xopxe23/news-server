package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xopxe23/news-server/internal/domain"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UsersRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	GetBookmarks(ctx context.Context, userId int) ([]domain.ArticleOutput, error)
}

type SessionsRepository interface {
	Create(ctx context.Context, token domain.RefreshSession) error
	GetToken(ctx context.Context, token string) (domain.RefreshSession, error)
}

type UsersService struct {
	repo         UsersRepository
	hasher       PasswordHasher
	sessionsRepo SessionsRepository
	hmacSecret   []byte
}

func NewUsersService(repo UsersRepository, hasher PasswordHasher, sessionsRepo SessionsRepository, secret []byte) *UsersService {
	return &UsersService{
		repo:         repo,
		sessionsRepo: sessionsRepo,
		hasher:       hasher,
		hmacSecret:   secret,
	}
}

func (s *UsersService) SignUp(ctx context.Context, input domain.SignUpInput) error {
	password, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}
	user := domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: password,
	}

	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Create(ctx, user)
}

func (s *UsersService) SignIn(ctx context.Context, input domain.SignInInput) (string, string, error) {
	password, err := s.hasher.Hash(input.Password)
	if err != nil {
		return "", "", err
	}

	user, err := s.repo.GetByCredentials(ctx, input.Email, password)
	if err != nil {
		return "", "", err
	}

	return s.generateTokens(ctx, user.Id)
}

func (s *UsersService) RefreshTokens(ctx context.Context, token string) (string, string, error) {
	session, err := s.sessionsRepo.GetToken(ctx, token)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", err
	}

	return s.generateTokens(ctx, session.UserId)
}

func (s *UsersService) generateTokens(ctx context.Context, userId int) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(userId),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	accessToken, err := t.SignedString(s.hmacSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := s.sessionsRepo.Create(ctx, domain.RefreshSession{
		UserId:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 12),
	}); err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (s *UsersService) ParseToken(ctx context.Context, token string) (int, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.hmacSecret, nil
	})
	if err != nil {
		return 0, err
	}
	if !t.Valid {
		return 0, errors.New("invalid token")
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}
	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}
	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}
	return id, nil
}

func (s *UsersService) GetBookmarks(ctx context.Context, userId int) ([]domain.ArticleOutput, error) {
	return s.repo.GetBookmarks(ctx, userId)
}
