package models

import (
	"github.com/ququzone/verifying-paymaster-service/db"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Address string `gorm:"type:varchar(42)"`
}

type ApiKeys struct {
	gorm.Model
	UserID      uint `json:"-"`
	User        User
	Key         string `gorm:"unique;type:varchar(30)"`
	Enable      bool
	Description string
}

func (a *ApiKeys) FindByKey(rep db.Repository, key string) (*ApiKeys, error) {
	var rec ApiKeys
	err := rep.Model(&ApiKeys{}).First(&rec, `"key" = ?`, key).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return &rec, err
}
