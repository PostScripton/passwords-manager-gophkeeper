package services

import (
	"context"
	pb "github.com/PostScripton/passwords-manager-gophkeeper/api/proto"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
)

type User interface {
	Register(ctx context.Context, login, password, aesSecret, privateKey string) error
	Login(ctx context.Context, login, password string) error
	Delete(ctx context.Context) error
}

type Auth interface {
	CheckAuthorized(ctx context.Context) (bool, error)
}

type SecureKeys interface {
	GenerateKeys() (string, string, error) // Возвращает зашифрованные AES и приватный ключи
	GetAesSecret(encAesSecret, encPrivateKey string) ([]byte, error)
}

type Services struct {
	User
	Auth
	SecureKeys
}

func NewServices(userClient pb.UserClient, repos *repository.Repository, jwtSecret, masterPassword string) *Services {
	return &Services{
		User:       NewUserService(userClient, repos.Settings),
		Auth:       NewAuthService(jwtSecret, repos.Settings),
		SecureKeys: NewSecureKeysService(masterPassword),
	}
}
