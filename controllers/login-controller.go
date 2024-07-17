package controllers

import (
	"net/http"

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

type loginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *LoginController) Login(ctx *gin.Context) {
	var form loginForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid form",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "login",
	})
}
