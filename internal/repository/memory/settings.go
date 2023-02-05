package memory

import (
	"context"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"sync"
)

type SettingsRepository struct {
	storage map[string]string
	mu      *sync.RWMutex
}

var _ repository.Settings = (*SettingsRepository)(nil)

func NewSettingsRepository() repository.Settings {
	return &SettingsRepository{
		storage: make(map[string]string),
		mu:      &sync.RWMutex{},
	}
}

func (r *SettingsRepository) Get(_ context.Context, key string) (string, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	value, existed := r.storage[key]

	return value, existed, nil
}

func (r *SettingsRepository) Set(_ context.Context, key, value string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, existed := r.storage[key]
	r.storage[key] = value

	return existed, nil
}

func (r *SettingsRepository) Delete(_ context.Context, key string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, existed := r.storage[key]
	delete(r.storage, key)

	return existed, nil
}

func (r *SettingsRepository) Truncate(_ context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage = make(map[string]string)

	return nil
}
