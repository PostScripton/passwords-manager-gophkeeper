package commands

import (
	"errors"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/cli/services"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:       "login",
	Short:     "logins you into the system",
	Args:      cobra.MinimumNArgs(2),
	ValidArgs: []string{"login", "password"},
	Run: func(cmd *cobra.Command, args []string) {
		authorized, err := dependencies.Services.Auth.CheckAuthorized(cmd.Context())
		if err != nil {
			authError(cmd, err)
			if err := dependencies.Services.User.Delete(cmd.Context()); err != nil {
				log.Error().Err(err).Msg("Deleting user info")
			}
			return
		}
		if authorized {
			cmd.PrintErrln("You are already logged into the system. " +
				"If you want to login via another user then logout first.")
			return
		}

		login, password := args[0], args[1]

		if err := dependencies.Services.User.Login(cmd.Context(), login, password); err != nil {
			if errors.Is(err, services.ErrCredentialsDontMatch) {
				cmd.PrintErrln("The credentials don't match any of our records")
				return
			}

			log.Error().Err(err).Msg("Logging user from command")
			return
		}

		cmd.Printf("[%s] Logged in successfully!\n", login)
	},
}
