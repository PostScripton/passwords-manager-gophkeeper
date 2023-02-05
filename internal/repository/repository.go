package repository

import (
	"context"
	"errors"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
)

var ErrLoginTaken = errors.New("this login is already taken")

type Users interface {
	FindByLogin(ctx context.Context, login string) (*models.User, error)
	Create(ctx context.Context, login, password, aesSecret, privateKey string) error
}

type Settings interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key, value string) (bool, error)
	Delete(ctx context.Context, key string) (bool, error)
	Truncate(ctx context.Context) error
}

type Repository struct {
	Users
	Settings
}

func NewRepository(factory Factory) *Repository {
	return &Repository{
		Users:    factory.CreateUserRepository(),
		Settings: factory.CreateSettingsRepository(),
	}
}
