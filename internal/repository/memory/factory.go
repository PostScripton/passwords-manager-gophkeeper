package memory

import "github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"

type Factory struct{}

var _ repository.Factory = (*Factory)(nil)

func NewFactory() repository.Factory {
	return &Factory{}
}

func (f *Factory) CreateUsersRepository() repository.Users {
	return NewUsersRepository()
}

func (f *Factory) CreateSettingsRepository() repository.Settings {
	return NewSettingsRepository()
}

func (f *Factory) CreateCredsSecretsRepository() repository.CredsSecrets {
	return NewCredsSecretsRepository()
}
