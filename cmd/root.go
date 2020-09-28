package cmd

import (
	"fmt"
	"github.com/psy-core/psysswd-vault/config"
	"os"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "A password vault for your password security.",
	Long:  `A password vault for your password security.`,
	Run:   runLogin,
}

func init() {
	rootCmd.PersistentFlags().StringP("conf", "c", "", "config file")
	rootCmd.PersistentFlags().StringP("username", "u", "", "give your username")
	rootCmd.PersistentFlags().StringP("password", "p", "", "give your master password")
}

func Execute() {

	checkError(rootCmd.Execute())

}

func checkError(e error) {
	if e != nil {
		fmt.Println("unexpect error: ", e)
		os.Exit(1)
	}
}

func readUsernameAndPassword(cmd *cobra.Command, conf *config.VaultConfig) (string, string, error) {
	//read username
	username, err := cmd.Flags().GetString("username")
	if err != nil {
		return "", "", err
	}

	if username == "" {
		username = conf.UserConf.DefaultUserName
	}

	if username == "" {
		fmt.Println("please give your master user name with -u or in config file.")
		os.Exit(1)
	}

	//read the master password
	password, err := cmd.Flags().GetString("password")
	if err != nil {
		return "", "", err
	}

	if password == "" {
		fmt.Print("please input your master password: ")
		passwordBytes, err := gopass.GetPasswdMasked()
		if err != nil {
			return "", "", err
		}
		password = string(passwordBytes)
	}

	return username, password, nil
}
