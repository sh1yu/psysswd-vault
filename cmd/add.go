package cmd

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/psy-core/psysswd-vault/config"
	"github.com/psy-core/psysswd-vault/persist"
	"github.com/spf13/cobra"
	"os"
)

var addCmd = &cobra.Command{
	Use:   "add <account-name> <account-user> [extra-message]",
	Short: "add a new account for given username",
	Long:  `add a new account info for given username`,
	Args:  cobra.RangeArgs(2, 3),
	Run:   func (cmd *cobra.Command, args []string) {
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
		runAdd(vaultConf.PersistConf.DataFile, username, password, args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(dataFile, username, password string, args []string) {
	account := args[0]
	user := args[1]
	extra := ""
	if len(args) == 3 {
		extra = args[2]
	}

	fmt.Printf("input password for account %s: ", account)
	passwordBytes, err := gopass.GetPasswdMasked()
	checkError(err)
	passwd := string(passwordBytes)

	saveData := &persist.DecodedRecord{
		Name:          account,
		Description:   "",
		LoginName:     user,
		LoginPassword: passwd,
		ExtraMessage:  extra,
	}
	err = persist.ModifyRecord(dataFile, username, password, saveData)
	checkError(err)

	fmt.Printf("add account %s success.\n", account)
}
