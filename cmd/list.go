package cmd

import (
	"fmt"
	"github.com/psy-core/psysswd-vault/config"
	"github.com/psy-core/psysswd-vault/persist"
	"github.com/spf13/cobra"
	"os"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list account info for given username",
	Long:  `list account info for given username`,
	Args:  cobra.NoArgs,
	Run:   func(cmd *cobra.Command, args []string) {
		vaultConf, err := config.InitConf(cmd.Flags().GetString("conf"))
		checkError(err)
		username, password, err := readUsernameAndPassword(cmd, vaultConf)
		checkError(err)

		exist, valid, err := persist.CheckUser(vaultConf.PersistConf.DataFile, username, password)
		checkError(err)
		if !exist {
			fmt.Println("user not registered: ", username)
			os.Exit(1)
		}
		if !valid {
			fmt.Println("Permission Denied.")
			os.Exit(1)
		}

		isPlain, err := cmd.Flags().GetBool("plain")
		checkError(err)

		runFind(isPlain, vaultConf.PersistConf.DataFile, username, password, "")
	},
}

func init() {
	listCmd.Flags().BoolP("plain", "P", false, "if true, print password in plain text")
	rootCmd.AddCommand(listCmd)
}

