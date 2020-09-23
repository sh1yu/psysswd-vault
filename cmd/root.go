package cmd

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "A password vault for your password security.",
	Long:  `A password vault for your password security.`,
	Run:   runLogin,
}

func init() {
	rootCmd.PersistentFlags().StringP("username", "u", "", "give your username")
	rootCmd.PersistentFlags().StringP("password", "p", "", "give your master password")
	err := rootCmd.MarkPersistentFlagRequired("username")
	checkError(err)
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

func readUsernameAndPassword(cmd *cobra.Command) (string, string, error) {
	//read username
	username, err := cmd.Flags().GetString("username")
	if err != nil {
		return "", "", nil
	}

	//read the master password
	password, err := cmd.Flags().GetString("password")
	if err != nil {
		return "", "", nil
	}

	if password == "" {
		fmt.Print("please input your master password: ")
		passwordBytes, err := gopass.GetPasswdMasked()
		if err != nil {
			return "", "", nil
		}
		password = string(passwordBytes)
	}

	return username, password, nil
}
