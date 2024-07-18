package database

type User struct {
	BaseModel
	Username string `json:"username"`
	Password string `json:"-"`
}

func InitModelUser() error {
	// return Db.AutoMigrate(&User{})
	return nil
}
