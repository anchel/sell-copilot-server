package controllers

import (
	"net/http"

	"github.com/anchel/sell-copilot-server/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	routes.AddRouteInitFunc(func(r *gin.RouterGroup) {
		ctl := NewExampleController()
		r.GET("/example/list", ctl.List)
	})
}

func NewExampleController() *ExampleController {
	return &ExampleController{
		BaseController: &BaseController{},
	}
}

type ExampleController struct {
	*BaseController
}

func (ctl *ExampleController) List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
	})
}
