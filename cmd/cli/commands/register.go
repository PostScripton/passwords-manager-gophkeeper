package commands

import (
	"errors"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/cli/services"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(registerCmd)
}

var registerCmd = &cobra.Command{
	Use:       "register",
	Short:     "registers you in the system",
	Args:      cobra.MinimumNArgs(2),
	ValidArgs: []string{"login", "password"},
	Run: func(cmd *cobra.Command, args []string) {
		authorized, err := dependencies.Services.Auth.CheckAuthorized(cmd.Context())
		if err != nil {
			authError(cmd, err)
			return
		}
		if authorized || errors.Is(err, services.ErrLoggedInAlready) {
			cmd.PrintErrln("You are already logged into the system. " +
				"If you want to register a new user then logout first.")
			return
		}

		login, password := args[0], args[1]

		cmd.Println("Generating secure keys...")
		aesSecret, privateKey, err := dependencies.Services.SecureKeys.GenerateKeys()
		if err != nil {
			log.Error().Err(err).Msg("Generating secure keys on registration command")
			cmd.PrintErrln("Error happened on generating secure keys :(")
			return
		}
		cmd.Println("The secure keys have been generated successfully!")

		if err := dependencies.Services.User.Register(cmd.Context(), login, password, aesSecret, privateKey); err != nil {
			if errors.Is(err, services.ErrLoginIsTaken) {
				cmd.PrintErrln("This login is already taken. Try another one.")
				return
			}

			log.Error().Err(err).Msg("Registering user from command")
			return
		}

		cmd.Printf("Registered [%s] successfully!\n", login)
	},
}
