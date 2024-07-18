package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anchel/sell-copilot-server/database"
	"github.com/anchel/sell-copilot-server/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func init() {
	routes.AddRouteInitFunc(func(r *gin.RouterGroup) {
		loginController := NewLoginController()
		r.POST("/login", loginController.Login)
		r.GET("/logout", loginController.Logout)
		r.GET("/userinfo", loginController.GetUserInfo)
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

func (c *LoginController) GetUserInfo(ctx *gin.Context) {
	session := sessions.Default(ctx)
	// log.Println("/api/userinfo", session.ID(), session.Get("login"))
	login, ok := session.Get("login").(int)
	if !ok || login != 1 {
		// log.Println("/api/userinfo", "no login")
		ctx.JSON(200, gin.H{
			"code":    1,
			"message": "no login",
		})
		return
	}
	// log.Println("/api/userinfo", "has login")
	userStr := session.Get("user").(string)
	var user database.User
	err := json.Unmarshal([]byte(userStr), &user)
	if err != nil {
		log.Println("/api/userinfo Unmarshal fail", err)
	}
	ctx.JSON(200, gin.H{
		"code":    0,
		"message": "ok",
		"user":    user,
	})
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
	var user database.User
	result := database.Db.Where(&database.User{Username: form.Username, Password: form.Password}).First(&user)
	if result.Error != nil {
		ctx.JSON(200, gin.H{
			"code":    1,
			"message": result.Error.Error(),
		})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(200, gin.H{
			"code":    1,
			"message": "user not found",
		})
		return
	}

	session := sessions.Default(ctx)
	session.Set("login", 1)
	userBS, err := json.Marshal(user)
	if err == nil {
		log.Println("/api/login marshal ok", string(userBS))
		session.Set("user", string(userBS))
	} else {
		log.Println("/api/login marshal fail", err)
	}

	if err := session.Save(); err != nil {
		ctx.JSON(200, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    0,
		"message": "ok",
		"user":    user,
	})
}

func (c *LoginController) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()

	if err := session.Save(); err != nil {
		ctx.JSON(200, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    0,
		"message": "ok",
	})
}
