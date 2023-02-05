package memory

import (
	"context"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"math/rand"
	"sync"
)

type UsersRepository struct {
	storage map[string]models.User
	mu      *sync.RWMutex
}

var _ repository.Users = (*UsersRepository)(nil)

func NewUsersRepository() repository.Users {
	return &UsersRepository{
		storage: make(map[string]models.User),
		mu:      &sync.RWMutex{},
	}
}

func (r *UsersRepository) FindByLogin(_ context.Context, login string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user, ok := r.storage[login]; ok {
		return &user, nil
	}

	return nil, nil
}

func (r *UsersRepository) Create(_ context.Context, login, password, aesSecret, privateKey string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[login]; ok {
		return repository.ErrLoginTaken
	}

	user := models.User{
		ID:         rand.Int(),
		Login:      login,
		Password:   password,
		AesSecret:  aesSecret,
		PrivateKey: privateKey,
	}
	r.storage[login] = user

	return nil
}
