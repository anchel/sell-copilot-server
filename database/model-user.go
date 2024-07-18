package database

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"-"`
}

func InitModelUser() error {
	// return Db.AutoMigrate(&User{})
	return nil
}
