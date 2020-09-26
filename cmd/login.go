package cmd

import (
	"fmt"
	"github.com/psy-core/psysswd-vault/config"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login vault and get a command shell",
	Long:  `login with master password, and get command shell`,
	Args:  cobra.NoArgs,
	Run:   runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) {
	vaultConf, err := config.InitConf(cmd.Flags().GetString("conf"))
	checkError(err)
	username, password, err := readUsernameAndPassword(cmd, vaultConf)
	checkError(err)
	fmt.Println("your username: ", username)
	fmt.Println("your password: ", password)

	//fixme check and give shell
}
