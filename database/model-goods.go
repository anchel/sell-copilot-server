package database

type Goods struct {
	BaseModel
	Title         string  `json:"title"`
	Description   *string `json:"description"`
	Thumbnail     *string `json:"thumbnail"`
	GoodsSkuTotal uint    `json:"goods_sku_total" gorm:"column:goods_sku_total"`
}

func InitModelGoods() error {
	// return Db.AutoMigrate(&Goods{})
	return nil
}
