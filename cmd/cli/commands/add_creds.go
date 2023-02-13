package commands

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCredsCmd)
}

var addCredsCmd = &cobra.Command{
	Use:       "add-creds",
	Short:     "saves your credentials from a web resource into the system",
	Args:      cobra.MinimumNArgs(3),
	ValidArgs: []string{"website", "login", "password", "additional_data"},
	Example:   `add-creds https://example.com/ qwerty 1234 {"additional_key": true}`, // TODO: accept flags instead of json
	Run: func(cmd *cobra.Command, args []string) {
		authorized, err := dependencies.Services.Auth.CheckAuthorized(cmd.Context())
		if err != nil {
			authError(cmd, err)
			if err := dependencies.Services.User.Delete(cmd.Context()); err != nil {
				log.Error().Err(err).Msg("Deleting user info")
			}
			return
		}
		if !authorized {
			unauthorized(cmd)
			return
		}
		id, err := dependencies.Services.Auth.GetID(cmd.Context())
		if err != nil {
			authError(cmd, err)
			return
		}

		additionalData := ""
		website, login, password := args[0], args[1], args[2]
		if len(args) > 3 {
			additionalData = args[3]
		}

		syncCreds(cmd)

		if err := dependencies.Services.CredsSecret.Add(cmd.Context(), id, website, login, password, additionalData); err != nil {
			log.Error().Err(err).Msg("Adding creds from command")
			return
		}

		uploadCreds(cmd)

		cmd.Printf("The credentials to the [%s] website have been added successfully!\n", website)
	},
}
