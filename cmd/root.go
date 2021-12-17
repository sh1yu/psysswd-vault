package cmd

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/psy-core/psysswd-vault/config"
	"github.com/psy-core/psysswd-vault/persist"
	"os"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "A password vault for your password security.",
	Long:  `A password vault for your password security.`,
	Run: func(cmd *cobra.Command, args []string) {
		vaultConf, username, password := runPreCheck(cmd)

		servePort, err := cmd.Flags().GetString("serve")
		checkError(err)

		runLogin(vaultConf.PersistConf.DataFile, username, password, servePort)
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("conf", "c", "", "config file")
	rootCmd.PersistentFlags().StringP("username", "u", "", "give your username")
	rootCmd.PersistentFlags().StringP("password", "p", "", "give your master password")
	rootCmd.PersistentFlags().StringP("serve", "s", "", "start a server for sync with given port")
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
