package services

import (
	"context"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
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

func (s *UserService) Create(
	ctx context.Context,
	login, hashedPassword, aesSecret, privateKey string,
) (*models.User, error) {
	user, err := s.repo.FindByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, repository.ErrLoginTaken
	}

	if err = s.repo.Create(ctx, login, hashedPassword, aesSecret, privateKey); err != nil {
		return nil, err
	}

	return s.repo.FindByLogin(ctx, login)
}

func (s *UserService) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	return s.repo.FindByLogin(ctx, login)
}
