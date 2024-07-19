package database

import (
	"log"
	"os"
	"time"

	libmysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{Logger: newLogger})
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
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
