package util

import (
	"crypto/rand"
	mathrand "math/rand"
)

const (
	saltMinLen = 8
	saltMaxLen = 32
)

func RandSalt() ([]byte, error) {
	// 生成8-32之间的随机数字
	salt := make([]byte, mathrand.Intn(saltMaxLen-saltMinLen)+saltMinLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}
