package services

import (
	"context"
	pb "github.com/PostScripton/passwords-manager-gophkeeper/api/proto"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
)

type User interface {
	Register(ctx context.Context, login, password, aesSecret, privateKey string) error
	Login(ctx context.Context, login, password string) error
	Delete(ctx context.Context) error
}

type Auth interface {
	CheckAuthorized(ctx context.Context) (bool, error)
	GetID(ctx context.Context) (int, error)
}

type SecureKeys interface {
	GenerateKeys() (string, string, error) // Возвращает зашифрованные AES и приватный ключи
	GetAesSecret(encAesSecret, encPrivateKey string) ([]byte, error)
	GetAesFromSettings(ctx context.Context) ([]byte, error)
}

type CredsSecret interface {
	Add(ctx context.Context, userID int, website, login, password, additionalData string) error
	Get(ctx context.Context, id int64) (*models.CredsSecret, error)
	Delete(ctx context.Context, id int64) error
	GetList(ctx context.Context, userID int) ([]*models.CredsSecret, error)
}

type Services struct {
	User
	Auth
	SecureKeys
	CredsSecret
}

func NewServices(
	userClient pb.UserClient,
	repos *repository.Repository,
	jwtSecret, masterPassword string,
) *Services {
	secureKeysService := NewSecureKeysService(masterPassword, repos.Settings)

	return &Services{
		User:        NewUserService(userClient, repos.Settings),
		Auth:        NewAuthService(jwtSecret, repos.Settings),
		SecureKeys:  secureKeysService,
		CredsSecret: NewCredsSecretService(repos.Settings, repos.CredsSecrets, secureKeysService),
	}
}
