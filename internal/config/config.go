package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig
	LoggerConfig
	ReposConfig
}

type ServerConfig struct {
	Address        string
	JWTSecret      string
	MasterPassword string
	SSLCertPath    string
	SSLKeyPath     string
}

type LoggerConfig struct {
	Level string
}

type ReposConfig struct {
	SQLite   string
	Postgres string
}

func NewConfig(configFolder string) *Config {
	viper.SetConfigType("yml")
	viper.AddConfigPath(configFolder)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Str("folder", configFolder).Msg("Reading in config")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal().Err(err).Msg("Unmarshalling config")
	}

	return &config
}
