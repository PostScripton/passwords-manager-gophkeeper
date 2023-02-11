package services

import (
	"context"
	"fmt"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	crypt "github.com/PostScripton/passwords-manager-gophkeeper/pkg/crypter"
)

type CredsSecretService struct {
	settingsRepo      repository.Settings
	credsSecretsRepo  repository.CredsSecrets
	secureKeysService SecureKeys
}

var _ CredsSecret = (*CredsSecretService)(nil)

func NewCredsSecretService(
	settingsRepo repository.Settings,
	credsSecretsRepo repository.CredsSecrets,
	secureKeys SecureKeys,
) CredsSecret {
	return &CredsSecretService{
		settingsRepo:      settingsRepo,
		credsSecretsRepo:  credsSecretsRepo,
		secureKeysService: secureKeys,
	}
}

func (s *CredsSecretService) Add(
	ctx context.Context,
	userID int,
	website, login, password, additionalData string,
) error {
	aesSecret, err := s.secureKeysService.GetAesFromSettings(ctx)
	if err != nil {
		return fmt.Errorf("get decrypted AES secret: %w", err)
	}

	encPassword, err := crypt.EncryptAES(aesSecret, password)
	if err != nil {
		return fmt.Errorf("encrypt password on adding credentials secret: %w", err)
	}

	if err = s.credsSecretsRepo.Create(ctx, userID, website, login, encPassword, additionalData); err != nil {
		return fmt.Errorf("store creds secret: %w", err)
	}

	return nil
}

func (s *CredsSecretService) Get(ctx context.Context, id int64) (*models.CredsSecret, error) {
	secret, err := s.credsSecretsRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get creds secret by id: %w", err)
	}
	if secret == nil {
		return nil, nil
	}

	aesSecret, err := s.secureKeysService.GetAesFromSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("get decrypted AES secret: %w", err)
	}

	password, err := crypt.DecryptAES(aesSecret, secret.Password)
	if err != nil {
		return nil, fmt.Errorf("decrypt password on getting credentials secret: %w", err)
	}

	secret.Password = password

	return secret, nil
}

func (s *CredsSecretService) Delete(ctx context.Context, id int64) error {
	if err := s.credsSecretsRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete creds secret: %w", err)
	}

	return nil
}

func (s *CredsSecretService) GetList(ctx context.Context, userID int) ([]*models.CredsSecret, error) {
	secrets, err := s.credsSecretsRepo.GetList(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get creds secret by id: %w", err)
	}

	return secrets, nil
}
