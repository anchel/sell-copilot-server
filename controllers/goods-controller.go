package controllers

import (
	"log"
	"net/http"

	"github.com/anchel/sell-copilot-server/database"
	"github.com/anchel/sell-copilot-server/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	routes.AddRouteInitFunc(func(r *gin.RouterGroup) {
		goodsController := NewGoodsController()
		r.GET("/goods/list", goodsController.List)
		r.POST("/goods/add", goodsController.Add)
		r.POST("/goods/del", goodsController.Del)
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

func (c *GoodsController) List(ctx *gin.Context) {
	var form listForm
	if err := ctx.ShouldBindQuery(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
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
		ctx.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "list",
		"list":    goods,
	})
}

func (c *GoodsController) Add(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "add",
	})
}

func (c *GoodsController) Del(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "del",
	})
}
