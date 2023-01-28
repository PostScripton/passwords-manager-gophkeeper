package commands

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"passwords-manager-gophkeeper/internal/config"
)

var (
	cfgFolder string
	cfg       *config.Config

	rootCmd = &cobra.Command{
		Use:   "pm",
		Short: "Passwords Manager GophKeeper",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Root command execution")
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFolder, "config", "c", "./configs", "path to config folder")
}

func initConfig() {
	cfg = config.NewConfig(cfgFolder)
}
