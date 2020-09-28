package cmd

import (
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/psy-core/psysswd-vault/config"
	"github.com/psy-core/psysswd-vault/internal/constant"
	"github.com/psy-core/psysswd-vault/persist"
	"github.com/spf13/cobra"
	"os"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list account info for given username",
	Long:  `list account info for given username`,
	Args:  cobra.NoArgs,
	Run:   runList,
}

func init() {
	listCmd.Flags().BoolP("plain", "P", false, "if true, print password in plain text")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	vaultConf, err := config.InitConf(cmd.Flags().GetString("conf"))
	checkError(err)
	username, password, err := readUsernameAndPassword(cmd, vaultConf)
	checkError(err)

	exist, valid, err := persist.CheckUser(vaultConf, username, password)
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

	printHeader := []string{"账号", "用户名", "密码", "额外信息", "更新时间"}
	printData := make([][]string, 0)

	decodeRecords, err := persist.QueryRecord(vaultConf, username, password, "")
	checkError(err)

	for _, record := range decodeRecords {
		if isPlain {
			printData = append(printData, []string{
				record.Name,
				record.LoginName,
				record.LoginPassword,
				record.ExtraMessage,
				record.UpdateTime.Format(constant.DateFormatSeconds),
			})
		} else {
			printData = append(printData, []string{
				record.Name,
				record.LoginName,
				string(bytes.Repeat([]byte("*"), len(record.LoginPassword))),
				record.ExtraMessage,
				record.UpdateTime.Format(constant.DateFormatSeconds),
			})
		}
	}

	tablePrint(printData, printHeader)
}

func tablePrint(data [][]string, header []string) {

	if len(data) == 0 {
		fmt.Println("+-----------------------+")
		fmt.Println("|      查询内容为空     |")
		fmt.Println("+-----------------------+")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
