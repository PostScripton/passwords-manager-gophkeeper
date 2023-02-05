package commands

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logoutCmd)
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logs you out of the system",
	Run: func(cmd *cobra.Command, args []string) {
		authorized, err := dependencies.Services.Auth.CheckAuthorized(cmd.Context())
		if err != nil {
			authError(cmd, err)
			return
		}
		if !authorized {
			unauthorized(cmd)
			return
		}

		if err := dependencies.Services.User.Delete(cmd.Context()); err != nil {
			log.Error().Err(err).Msg("Logging user from command")
			return
		}

		cmd.Println("Now you logged out.")
	},
}
