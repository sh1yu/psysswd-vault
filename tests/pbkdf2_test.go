package tests

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	mathrand "math/rand"
	"testing"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltMinLen = 8
	saltMaxLen = 32
	iter       = 1000
	keyLen     = 32
)

func TestPbkdf2(t *testing.T) {

	pwd := "This is a secret"

	salt, err := randSalt()
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

func randSalt() ([]byte, error) {
	// 生成8-32之间的随机数字
	salt := make([]byte, mathrand.Intn(saltMaxLen-saltMinLen)+saltMinLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func encryptPwdWithSalt(pwd, salt []byte) []byte {
	pwd = append(pwd, salt...)
	return pbkdf2.Key(pwd, salt, iter, keyLen, sha256.New)
}
