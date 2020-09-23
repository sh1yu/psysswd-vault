package cmd

import (
	"fmt"
	"os"

	"github.com/psy-core/psysswd-vault/internal/auth"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <account-name> <account-password>",
	Short: "add a new account for given username",
	Long:  `add a new account info for given username`,
	Args:  cobra.ExactArgs(2),
	Run:   runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) {
	username, password, err := readUsernameAndPassword(cmd)
	checkError(err)

	if !auth.Auth(username, password) {
		fmt.Println("Permission Denied.")
		os.Exit(1)
	}

	accountUser, accountPasswd := args[0], args[1]

	fmt.Println(accountUser, accountPasswd)
}
