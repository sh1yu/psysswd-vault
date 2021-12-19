package cmd

import (
	"fmt"

	"github.com/howeyc/gopass"
	"github.com/sh1yu/psysswd-vault/config"
	"github.com/sh1yu/psysswd-vault/persist"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register <master-account-name>",
	Short: "register a new master account for storage password",
	Long:  `register a new master account for storage password`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vaultConf, err := config.InitConf(cmd.Flags().GetString("conf"))
		checkError(err)
		err = runRegister(vaultConf.PersistConf.DataFile, args[0])
		checkError(err)
		fmt.Println("register success.")
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}

func runRegister(dataFile, accountUser string) error {
	fmt.Printf("Please input your new password for account '%s' :", accountUser)
	passwordBytes, err := gopass.GetPasswdMasked()
	if err != nil {
		return err
	}

	return persist.ModifyUser(dataFile, accountUser, string(passwordBytes))

}
