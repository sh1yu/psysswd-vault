package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"github.com/sh1yu/psysswd-vault/config"
	"github.com/sh1yu/psysswd-vault/internal/constant"
	"golang.org/x/crypto/pbkdf2"
	"io/ioutil"
)

//fixme deprecated
func ModifyAccount(conf *config.VaultConfig, username, password string) error {
	//metaFilePath, err := config.CreateFileIfNeeded("./meta.data")
	//if err != nil {
	//	return err
	//}

	content, _ := ioutil.ReadFile("./meta.data")

	var buf bytes.Buffer
	var userLen, passLen int32
	var offset int32 = 0

	changed := false

	for offset < int32(len(content)) {

		initOffset := offset

		err := binary.Read(bytes.NewBuffer(content[offset:4+offset]), binary.LittleEndian, &userLen)
		if err != nil {
			return err
		}
		offset += 4

		user := string(content[offset : userLen+offset])
		offset += userLen

		err = binary.Read(bytes.NewBuffer(content[offset:4+offset]), binary.LittleEndian, &passLen)
		if err != nil {
			return err
		}
		offset += 4

		pass := content[offset : passLen+offset]
		offset += passLen

		if user == username {
			rightPwd := pass[:32]
			salt := pass[32:]
			given := append([]byte(password), salt...)
			en := pbkdf2.Key(given, salt, constant.Pbkdf2Iter, 32, sha256.New)
			base64En := base64.StdEncoding.EncodeToString(en)
			if base64.StdEncoding.EncodeToString(rightPwd) == base64En {
				buf.Write(content[initOffset : offset-passLen-4])
				newPass := append(en, salt...)
				binary.Write(&buf, binary.LittleEndian, int32(len(newPass)))
				buf.Write(newPass)
			} else {
				buf.Write(content[initOffset:offset])
			}
			changed = true
		} else {
			buf.Write(content[initOffset:offset])
		}
	}

	if !changed {
		salt, err := RandSalt()
		if err != nil {
			return err
		}
		given := append([]byte(password), salt...)
		en := pbkdf2.Key(given, salt, constant.Pbkdf2Iter, 32, sha256.New)
		newPass := append(en, salt...)

		binary.Write(&buf, binary.LittleEndian, int32(len(username)))
		buf.Write([]byte(username))
		binary.Write(&buf, binary.LittleEndian, int32(len(newPass)))
		buf.Write(newPass)
	}

	return ioutil.WriteFile("./meta.data", buf.Bytes(), 0644)
}

//fixme deprecated
func Auth(conf *config.VaultConfig, username, password string) (bool, bool) {

	//metaFilePath, err := config.CreateFileIfNeeded("meta.data")
	//if err != nil {
	//	fmt.Println("some unexpected error happen for meta file:", err)
	//	os.Exit(1)
	//}
	content, _ := ioutil.ReadFile("./meta.data")

	var userLen, passLen int32
	var offset int32 = 0

	for offset < int32(len(content)) {
		err := binary.Read(bytes.NewBuffer(content[offset:4+offset]), binary.LittleEndian, &userLen)
		if err != nil {
			panic(err)
		}
		offset += 4

		user := string(content[offset : userLen+offset])
		offset += userLen

		err = binary.Read(bytes.NewBuffer(content[offset:4+offset]), binary.LittleEndian, &passLen)
		if err != nil {
			panic(err)
		}
		offset += 4

		pass := content[offset : passLen+offset]
		offset += passLen

		if user == username {
			rightPwd := pass[:32]
			salt := pass[32:]

			given := append([]byte(password), salt...)
			en := pbkdf2.Key(given, salt, constant.Pbkdf2Iter, 32, sha256.New)
			return true, base64.StdEncoding.EncodeToString(rightPwd) == base64.StdEncoding.EncodeToString(en)
		}
	}

	return false, false
}
