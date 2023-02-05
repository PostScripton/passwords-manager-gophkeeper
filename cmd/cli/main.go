package main

import (
	"context"
	"fmt"
	"github.com/PostScripton/passwords-manager-gophkeeper/cmd/cli/commands"
	servicesPkg "github.com/PostScripton/passwords-manager-gophkeeper/internal/cli/services"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/client"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/config"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

const notAssigned = "N/A"

var (
	buildVersion = notAssigned
	buildTime    = notAssigned
	buildCommit  = notAssigned
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "03:04:05PM"})

	ctx := context.Background()

	log.Info().Msg(fmt.Sprintf("Build version: %s", buildVersion))
	log.Info().Msg(fmt.Sprintf("Build date: %s", buildTime))
	log.Info().Msg(fmt.Sprintf("Build commit: %s\n", buildCommit))

	cfg := config.NewConfig("./config")

	db, err := sqlite.NewSQLite(ctx, cfg.ReposConfig.SQLite)
	if err != nil {
		log.Fatal().Err(err).Msg("Connecting to the SQLite database")
	}

	factory := sqlite.NewFactory(db)
	repos := repository.NewRepository(factory)

	userClient := client.NewUserClient(ctx, cfg.ServerConfig.Address)
	services := servicesPkg.NewServices(userClient, repos, cfg.ServerConfig.JWTSecret, cfg.ServerConfig.MasterPassword)

	commands.Execute(ctx, &commands.Dependencies{
		Services: services,
	})
}
