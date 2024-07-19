package controllers

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/anchel/sell-copilot-server/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	routes.AddRouteInitFunc(func(r *gin.RouterGroup) {
		ctl := NewImageController()
		r.POST("/image/upload", ctl.Upload)
	})
}

func NewImageController() *ImageController {
	return &ImageController{
		BaseController: &BaseController{},
	}
}

type ImageController struct {
	*BaseController
}

type BindFile struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (ctl *ImageController) Upload(c *gin.Context) {
	var bindFile BindFile

	// Bind file
	if err := c.ShouldBind(&bindFile); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    2,
			"message": err.Error(),
		})
		return
	}

	file := bindFile.File

	t := time.Now()
	filename := fmt.Sprintf("upload-%d%s", t.UnixMicro(), filepath.Ext(file.Filename))

	exePath, err := os.Executable()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    2,
			"message": err.Error(),
		})
		return
	}
	wd := filepath.Dir(exePath)
	dstFilePath := filepath.Join(wd, "image", filename)
	log.Println("upload filepath", dstFilePath)

	if err := c.SaveUploadedFile(file, dstFilePath); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    2,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "ok",
		"imagePath": filepath.Join(os.Getenv("SERVE_HOST"), "/upload-image", filename),
	})
}
