package controllers

import (
	"github.com/anchel/sell-copilot-server/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	routes.AddRouteInitFunc(func(r *gin.Engine) {
		loginController := NewLoginController()
		r.POST("/api/login", loginController.Login)
	})
}

func NewLoginController() *LoginController {
	return &LoginController{
		BaseController: &BaseController{},
	}
}

type LoginController struct {
	*BaseController
}

func (c *LoginController) Login(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "login",
	})
}
