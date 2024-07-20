package database

type GoodsSku struct {
	BaseModel
	GoodsId     uint    `json:"goods_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Imgurl      *string `json:"imgurl"`
	Thumbnail   *string `json:"thumbnail"`
	Total       *uint   `json:"num_total" gorm:"column:num_total"`
	Remain      *uint   `json:"num_remain" gorm:"column:num_remain"`
}

func InitModelGoodsSku() error {
	// return Db.AutoMigrate(&Goods{})
	return nil
}
