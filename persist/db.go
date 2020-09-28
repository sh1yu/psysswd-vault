package persist

import (
	"github.com/jinzhu/gorm"
	"github.com/psy-core/psysswd-vault/config"
)

func initialDB(dbFile string) (*gorm.DB, error) {

	db, err := config.InitDBFile(dbFile)
	if err != nil {
		return nil, err
	}

	if !db.HasTable(&User{}) {
		err = db.CreateTable(&User{}).Error
		if err != nil {
			return nil, err
		}
	}

	if !db.HasTable(&AccountRecord{}) {
		err = db.CreateTable(&AccountRecord{}).Error
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
