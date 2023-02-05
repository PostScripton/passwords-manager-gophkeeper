package services

import (
	"context"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.Users
}

var _ User = (*UserService)(nil)

func NewUserService(repo repository.Users) User {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(ctx context.Context, login, password, aesSecret, privateKey string) (*models.User, error) {
	user, err := s.repo.FindByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, repository.ErrLoginTaken
	}

	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return nil, err
	}

	if err = s.repo.Create(ctx, login, hashedPassword, aesSecret, privateKey); err != nil {
		return nil, err
	}

	return s.repo.FindByLogin(ctx, login)
}

func (s *UserService) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	return s.repo.FindByLogin(ctx, login)
}

func (s *UserService) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
