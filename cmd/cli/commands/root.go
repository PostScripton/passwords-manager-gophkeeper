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

func syncCreds(cmd *cobra.Command) {
	cmd.Println("Syncing credentials from the server...")
	if err := dependencies.Services.Sync.SyncCreds(cmd.Context()); err != nil {
		log.Error().Err(err).Msg("Syncing creds from the server")
		cmd.PrintErrln("Error happened on syncing credentials from the server")
	} else {
		cmd.Println("The credentials have been synced")
	}
}

func uploadCreds(cmd *cobra.Command) {
	cmd.Println("Uploading credentials from the server...")
	if err := dependencies.Services.Sync.UploadCreds(cmd.Context()); err != nil {
		log.Error().Err(err).Msg("Uploading creds to the server")
		cmd.PrintErrln("Error happened on uploading credentials from the server")
	} else {
		cmd.Println("The credentials have been uploaded")
	}
}
