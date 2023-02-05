package repository

type Factory interface {
	CreateUserRepository() Users
	CreateSettingsRepository() Settings
}
