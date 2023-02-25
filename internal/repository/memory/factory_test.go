package memory

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFactory_CreateCredsSecretsRepository(t *testing.T) {
	factory := &Factory{}
	repo := factory.CreateCredsSecretsRepository()

	require.NotNil(t, repo)
	assert.IsType(t, &CredsSecretsRepository{}, repo, "CreateCredsSecretsRepository()")
}

func TestFactory_CreateSettingsRepository(t *testing.T) {
	factory := &Factory{}
	repo := factory.CreateSettingsRepository()

	require.NotNil(t, repo)
	assert.IsType(t, &SettingsRepository{}, repo, "CreateSettingsRepository()")
}

func TestFactory_CreateUserRepository(t *testing.T) {
	factory := &Factory{}
	repo := factory.CreateUsersRepository()

	require.NotNil(t, repo)
	assert.IsType(t, &UsersRepository{}, repo, "CreateUsersRepository()")
}

func TestNewFactory(t *testing.T) {
	factory := NewFactory()
	require.NotNil(t, factory)
}
