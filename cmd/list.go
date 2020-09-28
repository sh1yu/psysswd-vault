package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/psy-core/psysswd-vault/config"
	"github.com/psy-core/psysswd-vault/internal/constant"
	"github.com/psy-core/psysswd-vault/internal/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/pbkdf2"
	"os"
	"strings"
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

	exist, valid := util.Auth(vaultConf, username, password)
	if !exist {
		fmt.Println("user not registered: ", username)
		os.Exit(1)
	}
	if !valid {
		fmt.Println("Permission Denied.")
		os.Exit(1)
	}

	printHeader := []string{"账号", "用户名", "密码", "额外信息"}
	printData := make([][]string, 0)

	err = util.RangePersistData(vaultConf, func(key, data []byte) {

		strKey := string(key)
		if !strings.HasPrefix(strKey, username) {
			return
		}

		var saltLen int32
		binary.Read(bytes.NewBuffer(data[:4]), binary.LittleEndian, &saltLen)
		salt := data[4 : 4+saltLen]

		enKey := pbkdf2.Key([]byte(password), salt, constant.Pbkdf2Iter, 32, sha256.New)
		plainBytes, err := util.AesDecrypt(data[4+saltLen:], enKey)
		checkError(err)

		var jsonData map[string]string
		err = json.Unmarshal(plainBytes, &jsonData)
		checkError(err)

		isPlain, err := cmd.Flags().GetBool("plain")
		checkError(err)

		if isPlain {
			printData = append(printData, []string{
				strings.TrimPrefix(strKey, username),
				jsonData["user"],
				jsonData["password"],
				jsonData["extra"],
			})
		} else {
			printData = append(printData, []string{
				strings.TrimPrefix(strKey, username),
				jsonData["user"],
				string(bytes.Repeat([]byte("*"), len(jsonData["password"]))),
				jsonData["extra"],
			})
		}

	})
	checkError(err)

	tablePrint(printData, printHeader)
}

func tablePrint(data [][]string, header []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
