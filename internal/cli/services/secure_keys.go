package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	crypt "github.com/PostScripton/passwords-manager-gophkeeper/pkg/crypter"
)

type SecureKeysService struct {
	masterPassword string
	settingsRepo   repository.Settings
}

func NewSecureKeysService(masterPassword string, settingsRepo repository.Settings) SecureKeys {
	return &SecureKeysService{
		masterPassword: masterPassword,
		settingsRepo:   settingsRepo,
	}
}

// GenerateKeys создаёт AES ключ и приватный ключ, который в свою очередь шифрует AES.
//
// Возвращает ключи в зашифрованном виде строкой, где первый ключ - AES, а второй - приватный.
func (s *SecureKeysService) GenerateKeys() (string, string, error) {
	aesSecret := make([]byte, 32)
	_, err := rand.Read(aesSecret)
	if err != nil {
		return "", "", fmt.Errorf("generating random key as AES secret: %w", err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", fmt.Errorf("generating RSA keys: %w", err)
	}

	encAesSecretBytes, err := crypt.EncryptRSA(&privateKey.PublicKey, aesSecret)
	if err != nil {
		return "", "", fmt.Errorf("encrypting AES via private key: %w", err)
	}
	encAesSecret := hex.EncodeToString(encAesSecretBytes)

	masterPassword, err := hex.DecodeString(s.masterPassword)
	if err != nil {
		return "", "", fmt.Errorf("decoding hexed master-password: %w", err)
	}

	encPrivateKey, err := crypt.EncryptAES(masterPassword, string(crypt.PrivateKeyToPemBytes(privateKey)))
	if err != nil {
		return "", "", fmt.Errorf("encrypting private key via AES using master-password: %w", err)
	}

	return encAesSecret, encPrivateKey, nil
}

func (s *SecureKeysService) GetAesSecret(encAesSecret, encPrivateKey string) ([]byte, error) {
	masterPassword, err := hex.DecodeString(s.masterPassword)
	if err != nil {
		return []byte{}, fmt.Errorf("decoding hexed master-password: %w", err)
	}

	privateKeyPem, err := crypt.DecryptAES(masterPassword, encPrivateKey)
	if err != nil {
		return []byte{}, fmt.Errorf("decrypting private key via aes: %w", err)
	}

	privateKey, err := crypt.PrivateKeyFromPemBytes([]byte(privateKeyPem))
	if err != nil {
		return []byte{}, fmt.Errorf("getting private key from pem bytes: %w", err)
	}

	encAesSecretBytes, err := hex.DecodeString(encAesSecret)
	if err != nil {
		return []byte{}, fmt.Errorf("decoding hex to : %w", err)
	}

	aesSecret, err := crypt.DecryptRSA(privateKey, encAesSecretBytes)
	if err != nil {
		return []byte{}, fmt.Errorf("decrypting AES secret via private key: %w", err)
	}

	return aesSecret, nil
}

func (s *SecureKeysService) GetAesFromSettings(ctx context.Context) ([]byte, error) {
	encAesSecret, _, err := s.settingsRepo.Get(ctx, "aes_secret")
	if err != nil {
		return []byte{}, fmt.Errorf("getting encrypted AES secret from settings: %w", err)
	}
	encPrivateKey, _, err := s.settingsRepo.Get(ctx, "private_key")
	if err != nil {
		return []byte{}, fmt.Errorf("getting encrypted private key from settings: %w", err)
	}

	aesSecret, err := s.GetAesSecret(encAesSecret, encPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("decrypt AES secret: %w", err)
	}

	return aesSecret, nil
}
