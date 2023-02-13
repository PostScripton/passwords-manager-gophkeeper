package commands

import (
	"fmt"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	rootCmd.AddCommand(listCredsCmd)
}

var listCredsCmd = &cobra.Command{
	Use:   "list-creds",
	Short: "list all credentials stored in the system",
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

		syncCreds(cmd)

		secrets, err := dependencies.Services.CredsSecret.GetList(cmd.Context(), id)
		if err != nil {
			log.Error().Err(err).Msg("Getting list of creds from command")
			return
		}

		displayAllCredsSecrets(cmd, secrets)
	},
}

func displayAllCredsSecrets(cmd *cobra.Command, secrets []*models.CredsSecret) {
	line := strings.Repeat("-", 10)
	cmd.Println(fmt.Sprintf("\n%s %s %s", line, "List of credentials", line))

	var website string
	for _, secret := range secrets {
		if website != secret.Website {
			website = secret.Website

			cmd.Printf("\nWebsite: %s\n", website)
		}

		cmd.Printf("-- ID: [%d], Login: [%s]\n", secret.ID, secret.Login)
	}
}
