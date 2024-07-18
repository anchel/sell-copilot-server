package database

type Goods struct {
	BaseModel
	Title         string `json:"title"`
	Description   string `json:"description"`
	Thumbnail     string `json:"thumbnail"`
	GoodsSkuTotal int32  `json:"goods_sku_total"`
}

func InitModelGoods() error {
	// return Db.AutoMigrate(&Goods{})
	return nil
}
