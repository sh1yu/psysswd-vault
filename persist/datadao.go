package persist

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/psy-core/psysswd-vault/internal/constant"
	"github.com/psy-core/psysswd-vault/internal/util"
	"golang.org/x/crypto/pbkdf2"
	"time"
)

func DumpRecord(dataFile string, masterUserName string) ([]*AccountRecord, error) {
	db, err := initialDB(dataFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var datas []*AccountRecord
	err = db.
		Where("user_name = ?", masterUserName).
		Find(&datas).Error
	return datas, err
}

func ImportRecord(dataFile string, records []*AccountRecord) error {
	db, err := initialDB(dataFile)
	if err != nil {
		return err
	}
	defer db.Close()

	totalCount := 0
	insertCount := 0
	updateCount := 0
	ignoreCount := 0
	errCount := 0

	for _, record := range records {
		totalCount++
		var exist AccountRecord
		err = db.Where("user_name = ?", record.UserName).Where("name=?", record.Name).First(&exist).Error
		if err == gorm.ErrRecordNotFound {
			data := AccountRecord{
				UserName:        record.UserName,
				Name:            record.Name,
				Description:     record.Description,
				LoginName:       record.LoginName,
				Salt:            record.Salt,
				LoginPasswordEn: record.LoginPasswordEn,
				ExtraMessage:    record.ExtraMessage,
				CreateTime:      record.CreateTime,
				UpdateTime:      record.UpdateTime,
			}
			err = db.Save(&data).Error
			if err != nil {
				fmt.Println("import account", record.Name, "err: ", err)
				errCount++
			} else {
				insertCount++
			}
		} else {
			//数据库中存在的记录较老，需要更新
			if exist.UpdateTime.Before(record.UpdateTime) {
				exist.Description = record.Description
				exist.LoginName = record.LoginName
				exist.Salt = record.Salt
				exist.LoginPasswordEn = record.LoginPasswordEn
				exist.ExtraMessage = record.ExtraMessage
				exist.CreateTime = record.CreateTime
				exist.UpdateTime = record.UpdateTime
				err = db.Save(&exist).Error
				if err != nil {
					fmt.Println("import account", record.Name, "err: ", err)
					errCount++
				} else {
					updateCount++
				}
			} else {
				ignoreCount++
			}
		}
	}

	fmt.Printf("import complete. total: %d, insert: %d, update: %d, ignore:%d, err: %d\n",
		totalCount, insertCount, updateCount, ignoreCount, errCount)

	return nil
}

func QueryRecord(dataFile string, masterUserName, masterPassword string, recordNameKeyword string) ([]*DecodedRecord, error) {

	db, err := initialDB(dataFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var datas []*AccountRecord
	if recordNameKeyword == "" {
		err = db.
			Where("user_name = ?", masterUserName).
			Find(&datas).Error
	} else {
		err = db.
			Where("user_name = ?", masterUserName).
			Where("name like ?", "%"+recordNameKeyword+"%").
			Find(&datas).Error
	}

	resultRecord := make([]*DecodedRecord, 0, len(datas))
	for _, data := range datas {

		saltBytes, err := base64.StdEncoding.DecodeString(data.Salt)
		if err != nil {
			return nil, err
		}
		loginPasswordEnBytes, err := base64.StdEncoding.DecodeString(data.LoginPasswordEn)
		if err != nil {
			return nil, err
		}
		enKey := pbkdf2.Key([]byte(masterPassword), saltBytes, constant.Pbkdf2Iter, 32, sha256.New)
		plainBytes, err := util.AesDecrypt(loginPasswordEnBytes, enKey)
		if err != nil {
			return nil, err
		}

		resultRecord = append(resultRecord, &DecodedRecord{
			Name:          data.Name,
			Description:   data.Description,
			LoginName:     data.LoginName,
			LoginPassword: string(plainBytes),
			ExtraMessage:  data.ExtraMessage,
			CreateTime:    data.CreateTime,
			UpdateTime:    data.UpdateTime,
		})
	}
	return resultRecord, err
}

func ModifyRecord(dbFile, masterUserName, masterPassword string, newData *DecodedRecord) error {

	db, err := initialDB(dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	var oldData AccountRecord
	err = db.
		Where("user_name = ?", masterUserName).
		Where("name=?", newData.Name).
		First(&oldData).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {

		salt, err := util.RandSalt()
		if err != nil {
			return err
		}
		keyEn := pbkdf2.Key([]byte(masterPassword), salt, constant.Pbkdf2Iter, 32, sha256.New)
		encrypted, err := util.AesEncrypt([]byte(newData.LoginPassword), keyEn)
		if err != nil {
			return err
		}

		saveData := AccountRecord{
			UserName:        masterUserName,
			Name:            newData.Name,
			Description:     newData.Description,
			LoginName:       newData.LoginName,
			Salt:            base64.StdEncoding.EncodeToString(salt),
			LoginPasswordEn: base64.StdEncoding.EncodeToString(encrypted),
			ExtraMessage:    newData.ExtraMessage,
			CreateTime:      time.Now(),
			UpdateTime:      time.Now(),
		}
		return db.Save(&saveData).Error
	}

	saltBytes, err := base64.StdEncoding.DecodeString(oldData.Salt)
	if err != nil {
		return err
	}
	keyEn := pbkdf2.Key([]byte(masterPassword), saltBytes, constant.Pbkdf2Iter, 32, sha256.New)
	encrypted, err := util.AesEncrypt([]byte(newData.LoginPassword), keyEn)
	if err != nil {
		return err
	}

	oldData.Description = newData.Description
	oldData.LoginName = newData.LoginName
	oldData.LoginPasswordEn = base64.StdEncoding.EncodeToString(encrypted)
	oldData.ExtraMessage = newData.ExtraMessage
	oldData.UpdateTime = time.Now()
	return db.Save(&oldData).Error
}
