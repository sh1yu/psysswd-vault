package cmd

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/psy-core/psysswd-vault/config"
	"github.com/psy-core/psysswd-vault/persist"
	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign <master-account-name>",
	Short: "add a new master account for storage password",
	Long:  `add a new master account for storage password`,
	Args:  cobra.ExactArgs(1),
	Run:   runSign,
}

func init() {
	rootCmd.AddCommand(signCmd)
}

func runSign(cmd *cobra.Command, args []string) {
	vaultConf, err := config.InitConf(cmd.Flags().GetString("conf"))
	checkError(err)

	fmt.Printf("Please input your new password for account '%s' :", args[0])
	passwordBytes, err := gopass.GetPasswdMasked()
	checkError(err)

	err = persist.ModifyUser(vaultConf, args[0], string(passwordBytes))
	checkError(err)

	fmt.Println("sign success.")
}
