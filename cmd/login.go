package cmd

import (
	"fmt"

	"github.com/jaczerob/clerk-gopher/internal/sys"
	"github.com/jaczerob/clerk-gopher/internal/toontown"
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
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if loginPassword == "" {
			loginPassword, err = keyring.Get(KEYRING_SERVICE, loginUsername)
			cobra.CheckErr(err)
		}

		gameserver, playcookie, err := toontown.Login(loginUsername, loginPassword)
		cobra.CheckErr(err)

		dir, err := sys.GetDirectory()
		cobra.CheckErr(err)

		executable := fmt.Sprintf("%s/%s", dir, sys.GetExecutable())

		log.WithField("username", loginUsername).Info("logging into toontown, have fun!")

		err = sys.RunExecutable(executable, gameserver, playcookie, doPipe)
		cobra.CheckErr(err)
	},
}

func init() {
	loginCmd.Flags().StringVarP(&loginUsername, "username", "u", "", "your TTR username")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "your TTR password")
	loginCmd.Flags().BoolVar(&doPipe, "pipe", false, "whether to pipe stdout to clerk-gopher process")
	loginCmd.MarkFlagRequired("username")
	rootCmd.AddCommand(loginCmd)
}
