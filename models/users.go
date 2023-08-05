package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Address string `gorm:"type:varchar(42)"`
}

type ApiKeys struct {
	gorm.Model
	UserID      uint
	User        User
	Key         string `gorm:"unique;type:varchar(30)"`
	Enable      bool
	Description string
}
