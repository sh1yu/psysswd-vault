package tests

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/psy-core/psysswd-vault/internal/util"
	"testing"

	"golang.org/x/crypto/pbkdf2"
)

const (
	iter   = 1000
	keyLen = 32
)

func TestPbkdf2(t *testing.T) {

	pwd := "This is a secret"

	salt, err := util.RandSalt()
	if err != nil {
		return
	}

	en := encryptPwdWithSalt([]byte(pwd), salt)
	en = append(en, salt...)

	// var buf bytes.Buffer
	// binary.Write(&buf, binary.LittleEndian, int32(3))
	// buf.Write([]byte("psy"))
	// binary.Write(&buf, binary.LittleEndian, int32(len(en)))
	// buf.Write(en)
	// ioutil.WriteFile("1.data", buf.Bytes(), 0644)

	encrypt := base64.StdEncoding.EncodeToString(en)
	fmt.Println(encrypt)

}

func encryptPwdWithSalt(pwd, salt []byte) []byte {
	pwd = append(pwd, salt...)
	return pbkdf2.Key(pwd, salt, iter, keyLen, sha256.New)
}
