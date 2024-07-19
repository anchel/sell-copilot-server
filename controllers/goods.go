package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/anchel/sell-copilot-server/database"
	"github.com/anchel/sell-copilot-server/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	routes.AddRouteInitFunc(func(r *gin.RouterGroup) {
		ctl := NewGoodsController()
		r.GET("/goods/list", ctl.List)
		r.POST("/goods/add", ctl.Add)
		r.POST("/goods/:id", ctl.Edit)
		r.POST("/goods/del", ctl.Del)
	})
}

func NewGoodsController() *GoodsController {
	return &GoodsController{
		BaseController: &BaseController{},
	}
}

type GoodsController struct {
	*BaseController
}

type listForm struct {
	Offset int32 `form:"offset"`
	Limit  int32 `form:"limit"`
}

func (ctl *GoodsController) List(c *gin.Context) {
	var form listForm
	if err := c.ShouldBindQuery(&form); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "invalid form: " + err.Error(),
		})
		return
	}
	log.Println("/api/goods/list form", form)

	if form.Offset < 0 || form.Offset > 100000000 {
		form.Offset = 0
	}
	if form.Limit <= 0 || form.Limit > 1000 {
		form.Limit = 20
	}
	log.Println("/api/goods/list form", form)

	var goods []database.Goods
	result := database.Db.Limit(int(form.Limit)).Offset(int(form.Offset)).Find(&goods)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"list":    goods,
	})
}

/**
 * add a goods record
 *
 */
type addForm struct {
	Title         string  `json:"title" binding:"required"`
	Description   *string `json:"description"`
	GoodsSkuTotal uint    `json:"goods_sku_total"`
}

func (ctl *GoodsController) Add(c *gin.Context) {
	var form addForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "invalid form: " + err.Error(),
		})
		return
	}
	goods := database.Goods{
		Title:         form.Title,
		Description:   form.Description,
		GoodsSkuTotal: form.GoodsSkuTotal,
	}

	result := database.Db.Create(&goods)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"goods": map[string]any{
			"id":              goods.ID,
			"title":           goods.Title,
			"description":     goods.Description,
			"thumbnail":       goods.Thumbnail,
			"goods_sku_total": goods.GoodsSkuTotal,
			"created_at":      goods.CreatedAt,
			"updated_at":      goods.UpdatedAt,
		},
	})
}

/**
 * modify a goods record
 *
 */
type editForm struct {
	Title         string  `json:"title" binding:"required"`
	Description   *string `json:"description"`
	GoodsSkuTotal uint    `json:"goods_sku_total"`
}

func (ctl *GoodsController) Edit(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "id invalid",
		})
		return
	}

	var form editForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "invalid form: " + err.Error(),
		})
		return
	}

	findGoods := &database.Goods{}
	result := database.Db.Find(&findGoods, id)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": result.Error.Error(),
		})
		return
	}
	if result.RowsAffected <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "record not exists",
		})
		return
	}

	findGoods.Title = form.Title
	findGoods.Description = form.Description
	findGoods.GoodsSkuTotal = form.GoodsSkuTotal

	result = database.Db.Save(&findGoods)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"goods": map[string]any{
			"id":              findGoods.ID,
			"title":           findGoods.Title,
			"description":     findGoods.Description,
			"thumbnail":       findGoods.Thumbnail,
			"goods_sku_total": findGoods.GoodsSkuTotal,
			"created_at":      findGoods.CreatedAt,
			"updated_at":      findGoods.UpdatedAt,
		},
	})
}

/**
 * delete a goods record
 *
 */
type delForm struct {
	ID uint `json:"id" binding:"required"`
}

func (ctl *GoodsController) Del(c *gin.Context) {
	var form delForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "invalid form: " + err.Error(),
		})
		return
	}

	result := database.Db.Find(&database.Goods{}, form.ID)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": result.Error.Error(),
		})
		return
	}

	if result.RowsAffected <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "record not exists",
		})
		return
	}

	result = database.Db.Delete(&database.Goods{}, form.ID)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": result.Error.Error(),
		})
		return
	}

	if result.RowsAffected <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "record delete fail",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
	})
}