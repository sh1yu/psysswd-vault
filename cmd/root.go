package cmd

import (
	"errors"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sh1yu/psysswd-vault/config"
	"github.com/sh1yu/psysswd-vault/persist"
	"os"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "A password vault for your password security.",
	Long:  `A password vault for your password security.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
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
		fmt.Println("please give your master user name with -u or 'user.defaultUserName' in config file.")
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

func runPreCheck(cmd *cobra.Command) (*config.VaultConfig, string, string) {

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

	return vaultConf, username, password
}

func checkRemoteCredential(credentials []config.CredentialConfig, username, token string) error {
	for _, credential := range credentials {
		if username == credential.User && token == credential.Token {
			return nil
		}
	}

	return errors.New("credential invalid for user " + username)
}

func getRemoteCredential(credentials []config.CredentialConfig, username string) string {
	for _, credential := range credentials {
		if credential.User == username {
			return credential.Token
		}
	}

	return ""
}
