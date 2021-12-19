package persist

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sh1yu/psysswd-vault/internal/constant"
	"github.com/sh1yu/psysswd-vault/internal/util"
	"golang.org/x/crypto/pbkdf2"
	"time"
)

func ModifyUser(dataFile, username, password string) error {
	db, err := initialDB(dataFile)
	if err != nil {
		return err
	}
	defer db.Close()

	var user User
	err = db.Where("name=?", username).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		salt, err := util.RandSalt()
		if err != nil {
			return err
		}
		passToken := pbkdf2.Key([]byte(password), salt, constant.Pbkdf2Iter, 32, sha256.New)
		user = User{
			Name:        username,
			Description: "",
			Salt:        base64.StdEncoding.EncodeToString(salt),
			PassToken:   base64.StdEncoding.EncodeToString(passToken),
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		}

		return db.Save(&user).Error
	}

	//用户已存在
	saltBytes, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		return err
	}
	passToken := pbkdf2.Key([]byte(password), saltBytes, constant.Pbkdf2Iter, 32, sha256.New)
	user.PassToken = base64.StdEncoding.EncodeToString(passToken)

	return db.Save(&user).Error
}

func CheckUser(dataFile, username, password string) (bool, bool, error) {
	db, err := initialDB(dataFile)
	if err != nil {
		return false, false, err
	}
	defer db.Close()

	var user User
	err = db.Where("name=?", username).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, false, err
	}
	if err == gorm.ErrRecordNotFound {
		return false, false, nil
	}

	saltBytes, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		return true, false, err
	}
	given := pbkdf2.Key([]byte(password), saltBytes, constant.Pbkdf2Iter, 32, sha256.New)

	isMatch := base64.StdEncoding.EncodeToString(given) == user.PassToken
	return true, isMatch, nil
}
