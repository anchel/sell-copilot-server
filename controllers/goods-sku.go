package controllers

import (
	"net/http"

	"github.com/anchel/sell-copilot-server/database"
	"github.com/anchel/sell-copilot-server/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	routes.AddRouteInitFunc(func(r *gin.RouterGroup) {
		ctl := NewGoodsSkuController()
		r.GET("/goods/:goodsId/sku/list", ctl.List)
		r.POST("/goods/:goodsId/sku/add", ctl.Add)
		r.POST("/goods/:goodsId/sku/del", ctl.Del)
		r.POST("/goods/:goodsId/sku/:id", ctl.Edit)
	})
}

func NewGoodsSkuController() *GoodsSkuController {
	return &GoodsSkuController{
		BaseController: &BaseController{},
	}
}

type GoodsSkuController struct {
	*BaseController
}

// 查询列表
func (ctl *GoodsSkuController) List(c *gin.Context) {
	gid, err := ctl.getParamGoodsId(c)
	if err != nil {
		return
	}

	goods, err := ctl.checkGoods(c, uint(gid))
	if err != nil {
		return
	}

	var skuList []database.GoodsSku
	result := database.Db.Limit(1000).Offset(0).Where("goods_id", goods.ID).Find(&skuList)
	if result.Error != nil {
		ctl.returnFail(c, 1, result.Error.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"list":    skuList,
	})
}

// 添加
type addSkuForm struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	Imgurl      *string `json:"imgurl"`
	Thumbnail   *string `json:"thumbnail"`
	Total       int     `json:"num_total"`
	Remain      int     `json:"num_remain"`
}

func (ctl *GoodsSkuController) Add(c *gin.Context) {
	gid, err := ctl.getParamGoodsId(c)
	if err != nil {
		return
	}

	var form addSkuForm
	if err := c.ShouldBindJSON(&form); err != nil {
		ctl.returnFail(c, 1, "invalid form: "+err.Error())
		return
	}

	goods, err := ctl.checkGoods(c, gid)
	if err != nil {
		return
	}

	sku := database.GoodsSku{
		GoodsId:     goods.ID,
		Title:       form.Title,
		Description: form.Description,
		Imgurl:      form.Imgurl,
		Thumbnail:   form.Thumbnail,
		Total:       form.Total,
		Remain:      form.Remain,
	}

	result := database.Db.Create(&sku)
	if result.Error != nil {
		ctl.returnFail(c, 1, result.Error.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"sku":     sku,
	})
}

// 修改
type editSkuForm struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Imgurl      *string `json:"imgurl"`
	Thumbnail   *string `json:"thumbnail"`
	Total       *int    `json:"num_total"`
	Remain      *int    `json:"num_remain"`
}

func (ctl *GoodsSkuController) Edit(c *gin.Context) {
	gid, err := ctl.getParamGoodsId(c)
	if err != nil {
		return
	}
	skuId, err := ctl.getParamId(c)
	if err != nil {
		return
	}

	// 检查form
	var form editSkuForm
	if err := c.ShouldBindJSON(&form); err != nil {
		ctl.returnFail(c, 1, "invalid form: "+err.Error())
		return
	}

	// 检查sku是否存在
	sku, err := ctl.checkGoodsSku(c, gid, skuId)
	if err != nil {
		return
	}

	// 执行修改
	result := database.Db.Model(&sku).Updates(form)
	if result.Error != nil {
		ctl.returnFail(c, 1, result.Error.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"sku":     sku,
	})
}

// 删除
type delSkuForm struct {
	ID uint `json:"id" binding:"required"`
}

func (ctl *GoodsSkuController) Del(c *gin.Context) {
	gid, err := ctl.getParamGoodsId(c)
	if err != nil {
		return
	}

	var form delSkuForm
	if err := c.ShouldBindJSON(&form); err != nil {
		ctl.returnFail(c, 1, "invalid form: "+err.Error())
		return
	}

	sku, err := ctl.checkGoodsSku(c, gid, form.ID) // 检查sku是否存在
	if err != nil {
		return
	}

	result := database.Db.Delete(&database.GoodsSku{}, sku.ID)
	if result.Error != nil {
		ctl.returnFail(c, 1, result.Error.Error())
		return
	}

	if result.RowsAffected <= 0 {
		ctl.returnFail(c, 1, "record delete fail")
		return
	}

	ctl.returnOk(c)
}
