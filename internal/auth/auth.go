package auth

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"io/ioutil"

	"github.com/psy-core/psysswd-vault/internal/constant"
	"golang.org/x/crypto/pbkdf2"
)

var passwdMap map[string][]byte

func init() {

	passwdMap = make(map[string][]byte)

	content, err := ioutil.ReadFile("1.data")
	if err != nil {
		panic(err)
	}

	var userLen, passLen int32
	var offset int32 = 0

	for offset < int32(len(content)) {
		err = binary.Read(bytes.NewBuffer(content[offset:4+offset]), binary.LittleEndian, &userLen)
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

		passwdMap[user] = pass
	}
}

func Auth(username, password string) bool {

	if pwd, ok := passwdMap[username]; ok {
		rightPwd := pwd[:32]
		salt := pwd[32:]

		given := append([]byte(password), salt...)
		en := pbkdf2.Key(given, salt, constant.Pbkdf2Iter, 32, sha256.New)
		return base64.StdEncoding.EncodeToString(rightPwd) == base64.StdEncoding.EncodeToString(en)
	}
	return false
}

