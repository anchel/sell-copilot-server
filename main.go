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
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sellcopilot",
	Short: "sell copilot",
	Long:  "sell copilot, server and frontend",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			fmt.Println("配置文件路径为空")
			os.Exit(1)
		}
		err := run(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

//go:embed sell-copilot/dist
var frontend embed.FS

func run(configPath string) error {
	fmt.Println("configPath", configPath)
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return err
	}
	r := gin.Default()
	r.SetTrustedProxies(nil)

	store, err := redis.NewStore(3, "tcp", os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"), []byte("secret666"))
	if err != nil {
		log.Println("Error redis.NewStore")
		return err
	}
	r.Use(sessions.Sessions("mysession", store))

	r.Use(static.Serve("/", static.EmbedFolder(frontend, "sell-copilot/dist")))
	r.Static("/upload-image", "./image")

	// enable single page application
	r.NoRoute(func(c *gin.Context) {
		fmt.Println("NoRoute", c.Request.URL.Path)
		c.FileFromFS("/sell-copilot/dist/template.html", http.FS(frontend))
	})

	err = database.InitDB()
	if err != nil {
		log.Println("Error database.InitDB")
		return err
	}

	exclude := []string{
		"/api/login",
		"/api/userinfo",
	}
	group := r.Group("/api", func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if lo.Contains(exclude, path) { // exclude
			// log.Println("check login path exclude", path)
			ctx.Next()
		} else {
			session := sessions.Default(ctx)
			login, ok := session.Get("login").(int)
			if !ok || login != 1 {
				// log.Println("check login false")
				ctx.JSON(http.StatusForbidden, gin.H{
					"code":    1,
					"message": "no login",
				})
				ctx.Abort()
			} else {
				// log.Println("check login success")
				ctx.Next()
			}
		}
	})
	routes.InitRoutes(group)
	return r.Run(os.Getenv("LISTEN_ADDR"))
}
