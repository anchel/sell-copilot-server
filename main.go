package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/anchel/sell-copilot-server/controllers" // Import all controllers
	"github.com/anchel/sell-copilot-server/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := gin.Default()
	r.SetTrustedProxies(nil)

	store, err := redis.NewStore(3, "tcp", os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"), []byte("secret"))
	if err != nil {
		log.Fatal("Error redis.NewStore")
	}
	r.Use(sessions.Sessions("mysession", store))
	r.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("count", 1)
		session.Save()
		c.Next()
	})

	routes.InitRoutes(r)
	r.Run(":8080")
}

func test2() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/incr", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count += 1
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})
	r.Run(":8080")
}

func test() {
	fmt.Println("Hello, World!")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("TOKEN", os.Getenv("TOKEN"))

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
