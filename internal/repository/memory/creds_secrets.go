package memory

import (
	"context"
	"fmt"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"math/rand"
	"sort"
	"sync"
)

type CredsSecretsRepository struct {
	storage map[int64]models.CredsSecret
	mu      *sync.RWMutex
}

var _ repository.CredsSecrets = (*CredsSecretsRepository)(nil)

func NewCredsSecretsRepository() repository.CredsSecrets {
	return &CredsSecretsRepository{
		storage: make(map[int64]models.CredsSecret),
		mu:      &sync.RWMutex{},
	}
}

func (r *CredsSecretsRepository) Create(
	_ context.Context,
	userID int,
	website, login, encPassword, additionalData string,
) error {
	exists, err := r.checkCredsSecretExists(userID, website, login)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("credentials for this website exist")
	}

	credsSecret := models.CredsSecret{
		ID:             rand.Int63(),
		Website:        website,
		Login:          login,
		Password:       encPassword,
		AdditionalData: additionalData,
		UserID:         userID,
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[credsSecret.ID] = credsSecret

	return nil
}

func (r *CredsSecretsRepository) GetById(_ context.Context, id int64) (*models.CredsSecret, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	creds, ok := r.storage[id]
	if !ok {
		return nil, fmt.Errorf("creds with such id [%d] is not found", id)
	}

	return &creds, nil
}

func (r *CredsSecretsRepository) Delete(_ context.Context, id int64) error {
	delete(r.storage, id)

	return nil
}

func (r *CredsSecretsRepository) GetList(_ context.Context, userID int) ([]*models.CredsSecret, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	list := make([]*models.CredsSecret, 0)

	for _, secret := range r.storage {
		if secret.UserID == userID {
			list = append(list, &secret)
		}
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].Website == list[j].Website {
			return list[i].Login < list[j].Login
		}

		return list[i].Website < list[j].Website
	})

	return list, nil
}

func (r *CredsSecretsRepository) checkCredsSecretExists(userID int, website, login string) (bool, error) {
	for _, secret := range r.storage {
		if secret.UserID == userID && secret.Website == website && secret.Login == login {
			return true, nil
		}
	}

	return false, nil
}
