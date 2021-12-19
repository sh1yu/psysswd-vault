package cmd

import (
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/sh1yu/psysswd-vault/internal/constant"
	"github.com/sh1yu/psysswd-vault/persist"
	"github.com/spf13/cobra"
	"os"
)

var findCmd = &cobra.Command{
	Use:   "find <account-keyword>",
	Short: "find given account info",
	Long:  `find given account info`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		vaultConf, username, password := runPreCheck(cmd)

		isPlain, err := cmd.Flags().GetBool("plain")
		checkError(err)

		err = runFind(isPlain, vaultConf.PersistConf.DataFile, username, password, args[0])
		checkError(err)
	},
}

func init() {
	findCmd.Flags().BoolP("plain", "P", false, "if true, print password in plain text")
	rootCmd.AddCommand(findCmd)
}

func runFind(isPlain bool, dataFile, username, password, searchKey string) error {

	printHeader := []string{"账号", "用户名", "密码", "额外信息", "更新时间"}
	printData := make([][]string, 0)

	decodeRecords, err := persist.QueryRecord(dataFile, username, password, searchKey)
	if err != nil {
		return err
	}

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
	return nil
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
