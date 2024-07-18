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

	// r.Use(func(c *gin.Context) {
	// 	session := sessions.Default(c)
	// 	session.Set("count", 1)
	// 	session.Save()
	// 	c.Next()
	// })

	r.Use(func(c *gin.Context) {
		log.Println("middleware 111")
	})
	r.Use(func(c *gin.Context) {
		log.Println("middleware 222")
	})

	routes.InitRoutes(r)
	r.Run(os.Getenv("LISTEN_ADDR"))
}
