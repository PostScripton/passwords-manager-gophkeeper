package services

import (
	"context"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"github.com/golang-jwt/jwt/v4"
)

type User interface {
	Create(ctx context.Context, login, password, aesSecret, privateKey string) (*models.User, error)
	FindByLogin(ctx context.Context, login string) (*models.User, error)
}

type Auth interface {
	GetSecret() string
	LoginByUser(user *models.User) (string, error)
	Login(ctx context.Context, login, password string) (string, error)
	ParseJWT(tokenString string) (*jwt.Token, error)
	GetIDFromJWT(token *jwt.Token) (int, error)
}

type Creds interface {
	GetList(ctx context.Context, userID int) ([]*models.CredsSecret, error)
	SetList(ctx context.Context, list []models.CredsSecret) error
}

type Services struct {
	User
	Auth
	Creds
}

func NewServices(repo *repository.Repository, JWTSecret string) *Services {
	return &Services{
		User:  NewUserService(repo.Users),
		Auth:  NewAuthService(repo.Users, JWTSecret),
		Creds: NewCredsService(repo.CredsSecrets),
	}
}
