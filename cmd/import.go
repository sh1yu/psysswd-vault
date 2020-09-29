package cmd

import (
	"bufio"
	"errors"
	"github.com/psy-core/psysswd-vault/persist"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
	"time"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import account info for given username",
	Long:  `import account info for given username`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		vaultConf, username, _ := runPreCheck(cmd)

		targetType, err := cmd.Flags().GetString("type")
		checkError(err)
		objName, err := cmd.Flags().GetString("obj")
		checkError(err)
		err = runImport(vaultConf.PersistConf.DataFile, targetType, objName, username)
		checkError(err)
	},
}

func init() {
	importCmd.Flags().StringP("type", "t", "text", "import account data type. [text | other]")
	importCmd.Flags().StringP("obj", "o", "export.data", "import account data name. [text | other]")
	rootCmd.AddCommand(importCmd)
}

func runImport(dataFile, targetType, objName, username string) error {

	switch targetType {
	case "text":
		file, err := os.Open(objName)
		if err != nil {
			return err
		}
		reader := bufio.NewReader(file)

		records := make([]*persist.AccountRecord, 0)
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			token := strings.FieldsFunc(strings.TrimSpace(line), func(r rune) bool {
				return r == 0x1f
			})

			if len(token) != 9 {
				continue
			}

			createTime, _ := time.Parse(time.RFC3339, strings.TrimSpace(token[7]))
			updateTime, _ := time.Parse(time.RFC3339, strings.TrimSpace(token[8]))
			accountRecord := persist.AccountRecord{
				UserName:        strings.TrimSpace(token[0]),
				Name:            strings.TrimSpace(token[1]),
				Description:     strings.TrimSpace(token[2]),
				LoginName:       strings.TrimSpace(token[3]),
				Salt:            strings.TrimSpace(token[4]),
				LoginPasswordEn: strings.TrimSpace(token[5]),
				ExtraMessage:    strings.TrimSpace(token[6]),
				CreateTime:      createTime,
				UpdateTime:      updateTime,
			}
			records = append(records, &accountRecord)
		}

		return persist.ImportRecord(dataFile, records)
	default:
		return errors.New("invalid target types")
	}
}
