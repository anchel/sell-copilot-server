package database

type GoodsSku struct {
	BaseModel
	GoodsId     uint    `json:"goods_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Imgurl      *string `json:"imgurl"`
	Thumbnail   *string `json:"thumbnail"`
	Total       int     `json:"num_total" gorm:"column:num_total"`
	Remain      int     `json:"num_remain" gorm:"column:num_remain"`
}

func InitModelGoodsSku() error {
	// return Db.AutoMigrate(&Goods{})
	return nil
}
