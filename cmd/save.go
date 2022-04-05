package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"

	log "github.com/sirupsen/logrus"
)

var (
	saveUsername string
	savePassword string
)

const KEYRING_SERVICE string = "TTR"

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Saves the given username and password to your system's keyring",
	Long: `Saving the login allows you to log in with just your username using the
	login command. Using a username already saved will override your previous login.`,
	Example: "clerk-gopher save -u username -p password",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = keyring.Set(KEYRING_SERVICE, saveUsername, savePassword)
		if err != nil {
			return
		}

		log.WithField("username", saveUsername).Info("login saved")
		return
	},
}

func init() {
	saveCmd.Flags().StringVarP(&saveUsername, "username", "u", "", "your TTR username")
	saveCmd.Flags().StringVarP(&savePassword, "password", "p", "", "your TTR password")
	saveCmd.MarkFlagRequired("username")
	saveCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(saveCmd)
}
