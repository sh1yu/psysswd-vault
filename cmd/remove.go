package cmd

import (
	"fmt"
	"github.com/psy-core/psysswd-vault/persist"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <account-name>",
	Short: "remove a account for given username",
	Long:  `remove a account info for given username`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		vaultConf, username, _ := runPreCheck(cmd)
		err := runRemove(vaultConf.PersistConf.DataFile, username, args)
		checkError(err)

		fmt.Printf("remove account %s success.\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func runRemove(dataFile, username string, args []string) error {
	return persist.RemoveRecord(dataFile, username, args[0])
}
