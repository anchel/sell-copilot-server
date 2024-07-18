package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/anchel/sell-copilot-server/controllers" // Import all controllers
	"github.com/anchel/sell-copilot-server/database"
	"github.com/anchel/sell-copilot-server/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/samber/lo"
)

//go:embed sell-copilot/dist
var frontend embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := gin.Default()
	r.SetTrustedProxies(nil)

	store, err := redis.NewStore(3, "tcp", os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"), []byte("secret666"))
	if err != nil {
		log.Fatal("Error redis.NewStore")
	}
	r.Use(sessions.Sessions("mysession", store))

	r.Use(static.Serve("/", static.EmbedFolder(frontend, "sell-copilot/dist")))
	// enable single page application
	r.NoRoute(func(c *gin.Context) {
		fmt.Println("NoRoute", c.Request.URL.Path)
		c.FileFromFS("/sell-copilot/dist/template.html", http.FS(frontend))
	})

	err = database.InitDB()
	if err != nil {
		log.Fatal("Error database.InitDB")
	}

	exclude := []string{
		"/api/login",
		"/api/userinfo",
	}
	group := r.Group("/api", func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if lo.Contains(exclude, path) { // exclude
			log.Println("check login path exclude", path)
			ctx.Next()
		} else {
			session := sessions.Default(ctx)
			login, ok := session.Get("login").(int)
			if !ok || login != 1 {
				log.Println("check login false")
				ctx.JSON(200, gin.H{
					"code":    1,
					"message": "no login",
				})
				ctx.Abort()
			} else {
				log.Println("check login success")
				ctx.Next()
			}
		}
	})
	routes.InitRoutes(group)
	r.Run(os.Getenv("LISTEN_ADDR"))
}
