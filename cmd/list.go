package cmd

import (
	"fmt"
	"os"

	"github.com/psy-core/psysswd-vault/internal/auth"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list account info for given username",
	Long:  `list account info for given username`,
	Args:  cobra.NoArgs,
	Run:   runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	username, password, err := readUsernameAndPassword(cmd)
	checkError(err)

	if !auth.Auth(username, password) {
		fmt.Println("Permission Denied.")
		os.Exit(1)
	}

	fmt.Println(username, ", you  can read your account now!")
}
