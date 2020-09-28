package persist

import "time"

type User struct {
	ID          int64     `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Name        string    `gorm:"size:255;index;not null;unique"`
	Description string    `gorm:"size:255"`
	Salt        []byte    `gorm:"size:255;not null;"`
	PassToken   []byte    `gorm:"size:255"`
	CreateTime  time.Time `gorm:"type:datetime"`
	UpdateTime  time.Time `gorm:"type:datetime"`
}

func (User) TableName() string {
	return "user"
}

type DecodedUser struct {
	Name        string
	Description string
	Password    string
	CreateTime  time.Time
	UpdateTime  time.Time
}

type AccountRecord struct {
	ID              int64     `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	UserName        string    `gorm:"size:255;index;not null"`
	Name            string    `gorm:"size:255;index;not null;unique"`
	Description     string    `gorm:"size:255"`
	LoginName       string    `gorm:"size:255;index"`
	Salt            []byte    `gorm:"size:255;not null;"`
	LoginPasswordEn []byte    `gorm:"size:255;index"`
	ExtraMessage    string    `gorm:"size:255;index"`
	CreateTime      time.Time `gorm:"type:datetime"`
	UpdateTime      time.Time `gorm:"type:datetime"`
}

func (AccountRecord) TableName() string {
	return "account_record"
}

type DecodedRecord struct {
	Name          string
	Description   string
	LoginName     string
	LoginPassword string
	ExtraMessage  string
	CreateTime    time.Time
	UpdateTime    time.Time
}
