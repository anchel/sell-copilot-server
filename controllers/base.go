package controllers

import (
	"errors"
	"net/http"

	"github.com/anchel/sell-copilot-server/database"
	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

func (ctl *BaseController) returnOk(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
	})
}

func (ctl *BaseController) returnFail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": msg,
	})
}

func (ctl *BaseController) checkGoods(c *gin.Context, goodsId uint) (*database.Goods, error) {
	goods := database.Goods{}
	result := database.Db.Find(&goods, goodsId)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": result.Error.Error(),
		})
		return nil, result.Error
	}

	if result.RowsAffected <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "record not exists",
		})
		return nil, errors.New("record not exists")
	}

	return &goods, nil
}

func (ctl *BaseController) checkGoodsSku(c *gin.Context, goodsId uint, skuId uint) (*database.GoodsSku, error) {
	sku := database.GoodsSku{}
	result := database.Db.Find(&sku, skuId)
	if result.Error != nil {
		ctl.returnFail(c, 1, result.Error.Error())
		return nil, result.Error
	}

	if result.RowsAffected <= 0 {
		ctl.returnFail(c, 1, "record not exists")
		return nil, errors.New("record not exists")
	}

	if sku.GoodsId != goodsId {
		ctl.returnFail(c, 1, "goods_id not equal")
		return nil, errors.New("goods_id not equal")
	}

	return &sku, nil
}
