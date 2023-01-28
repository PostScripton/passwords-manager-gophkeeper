package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"passwords-manager-gophkeeper/cmd/cli/commands"
)

const notAssigned = "N/A"

var (
	buildVersion string
	buildTime    string
	buildCommit  string
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "03:04:05PM"})

	if buildVersion == "" {
		buildVersion = notAssigned
	}
	if buildTime == "" {
		buildTime = notAssigned
	}
	if buildCommit == "" {
		buildCommit = notAssigned
	}

	commands.Execute()

	log.Info().Msg(fmt.Sprintf("Build version: %s", buildVersion))
	log.Info().Msg(fmt.Sprintf("Build date: %s", buildTime))
	log.Info().Msg(fmt.Sprintf("Build commit: %s", buildCommit))
}
