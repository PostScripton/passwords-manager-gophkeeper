package commands

import (
	"context"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/cli/services"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	dependencies *Dependencies

	rootCmd = &cobra.Command{
		Use:   "gophkeeper",
		Short: "Passwords Manager GophKeeper",
	}
)

type Dependencies struct {
	Services *services.Services
}

func Execute(ctx context.Context, deps *Dependencies) {
	dependencies = deps
	rootCmd.SetContext(ctx)

	if err := rootCmd.Execute(); err != nil {
		log.Debug().Err(err).Msg("Root command execution")
		rootCmd.PrintErrln(err)
	}
}

func unauthorized(cmd *cobra.Command) {
	cmd.PrintErrln("You are not authorized in order to perform this action")
}

func authError(cmd *cobra.Command, err error) {
	log.Error().Err(err).Msg("Checking whether authorized user")
	cmd.PrintErrln("Error happened on checking authorization")
}
