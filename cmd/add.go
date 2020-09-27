package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/howeyc/gopass"
	"github.com/psy-core/psysswd-vault/config"
	"github.com/psy-core/psysswd-vault/internal/constant"
	"github.com/psy-core/psysswd-vault/internal/util"
	"golang.org/x/crypto/pbkdf2"

	"github.com/psy-core/psysswd-vault/internal/auth"
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

	if !auth.Auth(username, password) {
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

	//使用master password加盐生成aes-256的key
	salt, err := util.RandSalt()
	checkError(err)
	keyEn := pbkdf2.Key([]byte(password), salt, constant.Pbkdf2Iter, 32, sha256.New)
	encrypted, err := util.AesEncrypt(dataBytes, keyEn)
	checkError(err)

	//buf 存入需要保存的加密数据
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, int32(len(salt)))
	buf.Write(salt)
	buf.Write(encrypted)

	bodyData, err := ioutil.ReadFile("body.data")
	checkError(err)
	bodyOffset := len(bodyData)
	bodyLen := buf.Len()
	bodyData = append(bodyData, buf.Bytes()...)

	//存储的key由master user和account共同组成
	storeKey := pbkdf2.Key([]byte(username+data["account"]), []byte{}, constant.Pbkdf2Iter, 8, sha256.New)
	indexData, err := ioutil.ReadFile("index.data")
	checkError(err)

	for i := 0; i < len(indexData); i += 32 {
		if base64.StdEncoding.EncodeToString(storeKey) == base64.StdEncoding.EncodeToString(indexData[i:i+8]) {
			//已经存在，改密码

			var updateIndexBuf bytes.Buffer
			binary.Write(&updateIndexBuf, binary.LittleEndian, int64(bodyOffset))
			binary.Write(&updateIndexBuf, binary.LittleEndian, int32(bodyLen))
			updateByte := updateIndexBuf.Bytes()

			for j := 0; j < 12; j++ {
				indexData[i+8+j] = updateByte[j]
			}

			ioutil.WriteFile("body.data", bodyData, 0644)
			ioutil.WriteFile("index.data", indexData, 0644)
			return
		}
	}

	//不存在，添加user和密码
	userByte := []byte(username + data["account"])
	keyOffset := len(bodyData)
	keyLen := len(userByte)
	bodyData = append(bodyData, userByte...)

	var addIndexBuf bytes.Buffer
	addIndexBuf.Write(storeKey)
	binary.Write(&addIndexBuf, binary.LittleEndian, int64(bodyOffset))
	binary.Write(&addIndexBuf, binary.LittleEndian, int32(bodyLen))
	binary.Write(&addIndexBuf, binary.LittleEndian, int64(keyOffset))
	binary.Write(&addIndexBuf, binary.LittleEndian, int32(keyLen))
	indexData = append(indexData, addIndexBuf.Bytes()...)

	ioutil.WriteFile("body.data", bodyData, 0644)
	ioutil.WriteFile("index.data", indexData, 0644)
}
