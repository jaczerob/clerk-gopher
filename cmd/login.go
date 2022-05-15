package cmd

import (
	"fmt"

	"github.com/jaczerob/clerk-gopher/internal/sys"
	"github.com/jaczerob/clerk-gopher/internal/toontown/login"
	"github.com/jaczerob/clerk-gopher/internal/toontown/update"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"

	log "github.com/sirupsen/logrus"
)

var (
	loginUsername string
	loginPassword string
	doPipe        bool
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Logs you into Toontown Rewritten with the information given",
	Long: `Providing a username only will search your system's base password manager 
	for logins associated with Toontown Rewritten and your username, otherwise providing 
	a password as well will log you straight in with that info`,
	Example: "clerk-gopher login -u username [-p password]",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if loginPassword == "" {
			loginPassword, err = keyring.Get(KEYRING_SERVICE, loginUsername)
			if err != nil {
				return
			}
		}

		updater, err := update.NewUpdateClient()
		if err != nil {
			return
		}

		err = updater.Update()
		if err != nil {
			return
		}

		login := login.NewLoginClient()

		gameserver, playcookie, err := login.Login(loginUsername, loginPassword)
		if err != nil {
			return
		}

		executable := fmt.Sprintf("%s/%s", updater.Directory, sys.GetExecutable())

		log.WithField("username", loginUsername).Printf("logging into toontown, have fun!")

		return sys.RunExecutable(executable, gameserver, playcookie, doPipe)
	},
}

func init() {
	loginCmd.Flags().StringVarP(&loginUsername, "username", "u", "", "your TTR username")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "your TTR password")
	loginCmd.Flags().BoolVar(&doPipe, "pipe", false, "whether to pipe stdout to clerk-gopher process")
	loginCmd.MarkFlagRequired("username")
	rootCmd.AddCommand(loginCmd)
}
