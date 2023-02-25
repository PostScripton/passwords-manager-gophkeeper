package repository

type Factory interface {
	CreateUsersRepository() Users
	CreateSettingsRepository() Settings
	CreateCredsSecretsRepository() CredsSecrets
}
