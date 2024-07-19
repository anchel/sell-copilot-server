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
		ctl := NewLoginController()
		r.POST("/login", ctl.Login)
		r.GET("/logout", ctl.Logout)
		r.GET("/userinfo", ctl.GetUserInfo)
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

func (ctl *LoginController) GetUserInfo(c *gin.Context) {
	session := sessions.Default(c)
	// log.Println("/api/userinfo", session.ID(), session.Get("login"))
	login, ok := session.Get("login").(int)
	if !ok || login != 1 {
		// log.Println("/api/userinfo", "no login")
		c.JSON(http.StatusOK, gin.H{
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
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"user":    user,
	})
}

type loginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (ctl *LoginController) Login(c *gin.Context) {
	var form loginForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "invalid form",
		})
		return
	}
	var user database.User
	result := database.Db.Where(&database.User{Username: form.Username, Password: form.Password}).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": result.Error.Error(),
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "user not found",
		})
		return
	}

	session := sessions.Default(c)
	session.Set("login", 1)
	userBS, err := json.Marshal(user)
	if err == nil {
		log.Println("/api/login marshal ok", string(userBS))
		session.Set("user", string(userBS))
	} else {
		log.Println("/api/login marshal fail", err)
	}

	if err := session.Save(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"user":    user,
	})
}

func (ctl *LoginController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()

	if err := session.Save(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
	})
}
