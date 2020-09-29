package cmd

import (
    "fmt"
    "github.com/howeyc/gopass"
    "github.com/psy-core/psysswd-vault/internal/util"
    "github.com/psy-core/psysswd-vault/persist"
    "github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
    Use:   "add <account-name> <account-user> [extra-message]",
    Short: "add a new account for given username",
    Long:  `add a new account info for given username`,
    Args:  cobra.RangeArgs(2, 3),
    Run: func(cmd *cobra.Command, args []string) {

        vaultConf, username, password := runPreCheck(cmd)

        isGenerate, err := cmd.Flags().GetBool("genpass")
        checkError(err)

        runAdd(vaultConf.PersistConf.DataFile, username, password, isGenerate, args)
    },
}

func init() {
    addCmd.Flags().BoolP("genpass", "g", false, "if true, generate a random password")
    rootCmd.AddCommand(addCmd)
}

func runAdd(dataFile, username, password string, isGenerate bool, args []string) {
    account := args[0]
    user := args[1]
    extra := ""
    if len(args) == 3 {
        extra = args[2]
    }

    var passwd string
    if isGenerate {
        newPasswd, err := util.GenPass("base3", 16)
        checkError(err)
        passwd = newPasswd
    } else {
        fmt.Printf("input password for account %s: ", account)
        passwordBytes, err := gopass.GetPasswdMasked()
        checkError(err)
        passwd = string(passwordBytes)
    }

    saveData := &persist.DecodedRecord{
        Name:          account,
        Description:   "",
        LoginName:     user,
        LoginPassword: passwd,
        ExtraMessage:  extra,
    }
    err := persist.ModifyRecord(dataFile, username, password, saveData)
    checkError(err)

    fmt.Printf("add account %s success.\n", account)
}
