package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sh1yu/psysswd-vault/persist"
	"github.com/spf13/cobra"
	"io/ioutil"
	"time"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export account info for given username",
	Long:  `export account info for given username`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		vaultConf, username, _ := runPreCheck(cmd)

		targetType, err := cmd.Flags().GetString("type")
		checkError(err)
		objName, err := cmd.Flags().GetString("obj")
		checkError(err)

		err = runExport(vaultConf.PersistConf.DataFile, targetType, objName, username)
		checkError(err)
		fmt.Println("export complete.")
	},
}

func init() {
	exportCmd.Flags().StringP("type", "t", "text", "export account data type. [text | other]")
	exportCmd.Flags().StringP("obj", "o", "export.data", "export account data name. [text | other]")
	rootCmd.AddCommand(exportCmd)
}

func runExport(dataFile, targetType, objName, username string) error {

	datas, err := persist.DumpRecord(dataFile, username)
	if err != nil {
		return err
	}

	switch targetType {
	case "text":
		var buf bytes.Buffer
		for _, data := range datas {
			line := fmt.Sprintf("%s \x1f %s \x1f %s \x1f %s \x1f %s \x1f %s \x1f %s \x1f %s \x1f %s\n",
				data.UserName, data.Name, data.Description, data.LoginName, data.Salt, data.LoginPasswordEn, data.ExtraMessage,
				data.CreateTime.Format(time.RFC3339), data.UpdateTime.Format(time.RFC3339))
			buf.WriteString(line)
		}

		return ioutil.WriteFile(objName, buf.Bytes(), 0644)
	default:
		return errors.New("invalid target types")
	}
}
