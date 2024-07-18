package database

import "gorm.io/gorm"

type Goods struct {
	gorm.Model
	Title         string `json:"title"`
	Description   string `json:"description"`
	Thumbnail     string `json:"thumbnail"`
	GoodsSkuTotal int32  `json:"goods_sku_total"`
}

func InitModelGoods() error {
	// return Db.AutoMigrate(&Goods{})
	return nil
}
