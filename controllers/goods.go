package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/anchel/sell-copilot-server/lib/image"
	"github.com/anchel/sell-copilot-server/lib/util"
	"github.com/samber/lo"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/anchel/sell-copilot-server/database"
	"github.com/anchel/sell-copilot-server/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	routes.AddRouteInitFunc(func(r *gin.RouterGroup) {
		ctl := NewGoodsController()
		r.GET("/goods/list", ctl.List)
		r.POST("/goods/add", ctl.Add)
		r.POST("/goods/:goodsId", ctl.Edit)
		r.POST("/goods/del", ctl.Del)
		r.POST("/goods/merge-image", ctl.mergeImage)
	})
}

func NewGoodsController() *GoodsController {
	return &GoodsController{
		BaseController: &BaseController{},
	}
}

type GoodsController struct {
	*BaseController
}

type listForm struct {
	Offset int32 `form:"offset"`
	Limit  int32 `form:"limit"`
}

func (ctl *GoodsController) List(c *gin.Context) {
	var form listForm
	if err := c.ShouldBindQuery(&form); err != nil {
		ctl.returnFail(c, 1, "invalid form: "+err.Error())
		return
	}

	if form.Offset < 0 || form.Offset > 100000000 {
		form.Offset = 0
	}
	if form.Limit <= 0 || form.Limit > 1000 {
		form.Limit = 20
	}

	var goodsList []database.Goods
	result := database.Db.Limit(int(form.Limit)).Offset(int(form.Offset)).Find(&goodsList)
	if result.Error != nil {
		ctl.returnFail(c, 1, result.Error.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"list":    goodsList,
	})
}

/**
 * add a goods record
 *
 */
type addForm struct {
	Title         string  `json:"title" binding:"required"`
	Description   *string `json:"description"`
	GoodsSkuTotal uint    `json:"goods_sku_total"`
}

func (ctl *GoodsController) Add(c *gin.Context) {
	var form addForm
	if err := c.ShouldBindJSON(&form); err != nil {
		ctl.returnFail(c, 1, "invalid form: "+err.Error())
		return
	}
	goods := database.Goods{
		Title:         form.Title,
		Description:   form.Description,
		GoodsSkuTotal: form.GoodsSkuTotal,
	}

	result := database.Db.Create(&goods)
	if result.Error != nil {
		ctl.returnFail(c, 1, result.Error.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"goods":   goods,
	})
}

/**
 * modify a goods record
 * 因为编辑的时候，提交的各个字段都允许为空，当为空时表示不改变当前字段的值。所以下面的editform每个字段都用指针，当指针等于nil时，表示前端没有传，那么就不改变当前值
 * 但是这样会不会有一个问题，如果就是想要某个数据库字段变为null呢？
 */
type editForm struct {
	Title         *string `json:"title"`
	Description   *string `json:"description"`
	GoodsSkuTotal *uint   `json:"goods_sku_total"`
}

func (ctl *GoodsController) Edit(c *gin.Context) {
	gid, err := ctl.getParamGoodsId(c)
	if err != nil {
		return
	}

	var form editForm
	if err := c.ShouldBindJSON(&form); err != nil {
		ctl.returnFail(c, 1, "invalid form: "+err.Error())
		return
	}

	goods, err := ctl.checkGoods(c, gid)
	if err != nil {
		return
	}

	formStr, _ := json.Marshal(&form)
	var updateForm database.Goods
	err = json.Unmarshal(formStr, &updateForm)
	if err != nil {
		ctl.returnFail(c, 1, "json unmarshal fail: "+err.Error())
		return
	}

	result := database.Db.Model(goods).Updates(updateForm)
	if result.Error != nil {
		ctl.returnFail(c, 1, result.Error.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"goods":   goods,
	})
}

/**
 * delete a goods record
 *
 */
type delForm struct {
	ID uint `json:"id" binding:"required"`
}

func (ctl *GoodsController) Del(c *gin.Context) {
	var form delForm
	if err := c.ShouldBindJSON(&form); err != nil {
		ctl.returnFail(c, 1, "invalid form: "+err.Error())
		return
	}

	goods, err := ctl.checkGoods(c, form.ID) // 检查商品是否存在
	if err != nil {
		return
	}

	result := database.Db.Delete(&database.Goods{}, goods.ID)
	if result.Error != nil {
		ctl.returnFail(c, 1, result.Error.Error())
		return
	}

	if result.RowsAffected <= 0 {
		ctl.returnFail(c, 1, "record delete fail")
		return
	}

	ctl.returnOk(c)
}

type mergeImageForm struct {
	ID uint `json:"id" binding:"required"`
}

func (ctl *GoodsController) mergeImage(c *gin.Context) {
	var form mergeImageForm
	if err := c.ShouldBindJSON(&form); err != nil {
		ctl.returnFail(c, 1, "invalid form: "+err.Error())
		return
	}

	exePwd, err := util.GetExePwd()
	if err != nil {
		log.Println(err)
		ctl.returnFail(c, 1, err.Error())
		return
	}

	goods, err := ctl.checkGoods(c, form.ID) // 检查商品是否存在
	if err != nil {
		return
	}

	var skuList []database.GoodsSku
	result := database.Db.Limit(1000).Offset(0).Where("goods_id", goods.ID).Find(&skuList)
	if result.Error != nil {
		ctl.returnFail(c, 1, result.Error.Error())
		return
	}

	pathList := lo.Map(skuList, func(sku database.GoodsSku, index int) string {
		return filepath.Join(exePwd, *sku.Imgurl)
	})

	log.Println("pathList", pathList)

	is, err := image.NewList(pathList)
	if err != nil {
		log.Println("merge image:", err)
		ctl.returnFail(c, 1, err.Error())
		return
	}

	dx := 3
	if len(pathList) > 12 {
		dx = 4
	}
	img, err := is.ApplyGridLayout(dx, 100)
	if err != nil {
		log.Println(err)
		ctl.returnFail(c, 1, err.Error())
		return
	}

	t := time.Now()
	filename := fmt.Sprintf("upload-%d%s", t.UnixMicro(), ".jpg")
	dstFilePath := filepath.Join(exePwd, "upload-image", filename)

	// 保存图片
	err = image.Save(img, dstFilePath)
	if err != nil {
		log.Println(err)
		ctl.returnFail(c, 1, err.Error())
		return
	}

	// 更新数据库记录
	result = database.Db.Model(goods).Updates(map[string]interface{}{"thumbnail": filepath.Join("upload-image", filename)})
	if result.Error != nil {
		ctl.returnFail(c, 1, result.Error.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      0,
		"message":   "ok",
		"imagePath": filepath.Join(os.Getenv("SERVE_HOST"), "/upload-image", filename),
	})
}
