package database

import (
	"os"
	"time"

	libmysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() error {
	cfg := libmysql.NewConfig()
	cfg.User = os.Getenv("MYSQL_USER")
	cfg.Passwd = os.Getenv("MYSQL_PASSWORD")
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv("MYSQL_ADDR")
	cfg.DBName = os.Getenv("MYSQL_DB")
	cfg.ParseTime = true
	cfg.Loc = time.Local

	db, err := gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{})
	if err != nil {
		return err
	}
	Db = db
	err = InitModelUser()
	if err != nil {
		return err
	}
	err = InitModelGoods()
	if err != nil {
		return err
	}
	return nil
}

type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
