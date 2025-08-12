package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/huntercenter1/backend-test/user-service/internal/auth"
	"github.com/huntercenter1/backend-test/user-service/internal/models"
	"github.com/huntercenter1/backend-test/user-service/internal/repo"
)

var defaultTimeout = 5 * time.Second

type UserService interface {
	Create(ctx context.Context, username, email, password string) (*models.User, error)
	Get(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, id, username, email, password string) (*models.User, error)
	Delete(ctx context.Context, id string) error
	Authenticate(ctx context.Context, username, password string) (string, error)
	Validate(ctx context.Context, id string) (bool, error)
}

type userService struct {
	repo repo.UserRepo
}

func New(repo repo.UserRepo) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, username, email, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	username = strings.TrimSpace(username)
	email = strings.TrimSpace(strings.ToLower(email))
	if username == "" || email == "" || password == "" {
		return nil, errors.New("missing fields")
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
	}
	return s.repo.Create(ctx, user)
}

func (s *userService) Get(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) Update(ctx context.Context, id, username, email, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(username) != "" {
		u.Username = strings.TrimSpace(username)
	}
	if strings.TrimSpace(email) != "" {
		u.Email = strings.TrimSpace(strings.ToLower(email))
	}
	if strings.TrimSpace(password) != "" {
		hash, err := auth.HashPassword(password)
		if err != nil {
			return nil, err
		}
		u.PasswordHash = hash
	}
	return s.repo.Update(ctx, u)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) Authenticate(ctx context.Context, username, password string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	u, err := s.repo.GetByUsername(ctx, strings.TrimSpace(username))
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if !auth.CheckPassword(u.PasswordHash, password) {
		return "", errors.New("invalid credentials")
	}
	return u.ID, nil
}

func (s *userService) Validate(ctx context.Context, id string) (bool, error) {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
