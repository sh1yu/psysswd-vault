package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list account info for given username",
	Long:  `list account info for given username`,
	Args:  cobra.NoArgs,
	Run:   func(cmd *cobra.Command, args []string) {

		vaultConf, username, password := runPreCheck(cmd)

		isPlain, err := cmd.Flags().GetBool("plain")
		checkError(err)

		runFind(isPlain, vaultConf.PersistConf.DataFile, username, password, "")
	},
}

func init() {
	listCmd.Flags().BoolP("plain", "P", false, "if true, print password in plain text")
	rootCmd.AddCommand(listCmd)
}

