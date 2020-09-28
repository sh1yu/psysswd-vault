package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	"github.com/howeyc/gopass"
	"github.com/psy-core/psysswd-vault/config"
	"github.com/psy-core/psysswd-vault/internal/constant"
	"github.com/psy-core/psysswd-vault/internal/util"
	"golang.org/x/crypto/pbkdf2"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <account-name> <account-user> [extra-message]",
	Short: "add a new account for given username",
	Long:  `add a new account info for given username`,
	Args:  cobra.RangeArgs(2, 3),
	Run:   runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) {
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

	data := map[string]string{
		"account": args[0],
		"user":    args[1],
		"extra":   "",
	}

	if len(args) == 3 {
		data["extra"] = args[2]
	}

	fmt.Printf("input password for account %s: ", data["account"])
	passwordBytes, err := gopass.GetPasswdMasked()
	checkError(err)
	data["password"] = string(passwordBytes)

	dataBytes, err := json.Marshal(data)
	checkError(err)

	salt, err := util.RandSalt()
	checkError(err)
	keyEn := pbkdf2.Key([]byte(password), salt, constant.Pbkdf2Iter, 32, sha256.New)
	encrypted, err := util.AesEncrypt(dataBytes, keyEn)
	checkError(err)

	//finalData 存入需要保存的加密数据
	var finalData bytes.Buffer
	binary.Write(&finalData, binary.LittleEndian, int32(len(salt)))
	finalData.Write(salt)
	finalData.Write(encrypted)

	//存储的key由master user和account共同组成
	err = util.ModifyData(vaultConf, []byte(username+data["account"]), finalData.Bytes())
	checkError(err)

	fmt.Printf("add account %s success.\n", data["account"])
}
