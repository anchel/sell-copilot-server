package controllers

import (
	"github.com/anchel/sell-copilot-server/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	routes.AddRouteInitFunc(func(r *gin.Engine) {
		goodsController := NewGoodsController()
		r.GET("/api/goods/list", goodsController.List)
		r.POST("/api/goods/add", goodsController.Add)
		r.POST("/api/goods/del", goodsController.Del)
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

func (c *GoodsController) List(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "list",
	})
}

func (c *GoodsController) Add(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "add",
	})
}

func (c *GoodsController) Del(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "del",
	})
}
