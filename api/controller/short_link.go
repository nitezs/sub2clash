package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"sub2clash/config"
	"sub2clash/model"
	"sub2clash/utils"
	"sub2clash/utils/database"
	"sub2clash/validator"
	"time"
)

func ShortLinkGenHandler(c *gin.Context) {
	// 从请求中获取参数
	var params validator.ShortLinkGenValidator
	if err := c.ShouldBind(&params); err != nil {
		c.String(400, "参数错误: "+err.Error())
	}
	// 生成hash
	hash := utils.RandomString(config.Default.ShortLinkLength)
	// 存入数据库
	var item model.ShortLink
	result := database.DB.Model(&model.ShortLink{}).Where("url = ?", params.Url).First(&item)
	if result.Error == nil {
		c.String(200, item.Hash)
		return
	} else {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.String(500, "数据库错误: "+result.Error.Error())
			return
		}
	}
	// 如果记录存在则重新生成hash，直到记录不存在
	result = database.DB.Model(&model.ShortLink{}).Where("hash = ?", hash).First(&item)
	for result.Error == nil {
		hash = utils.RandomString(config.Default.ShortLinkLength)
		result = database.DB.Model(&model.ShortLink{}).Where("hash = ?", hash).First(&item)
	}
	// 创建记录
	database.DB.FirstOrCreate(
		&model.ShortLink{
			Hash:            hash,
			Url:             params.Url,
			LastRequestTime: -1,
		},
	)
	// 返回短链接
	c.String(200, hash)
}

func ShortLinkGetHandler(c *gin.Context) {
	// 获取动态路由
	hash := c.Param("hash")
	// 查询数据库
	var shortLink model.ShortLink
	result := database.DB.Where("hash = ?", hash).First(&shortLink)
	// 重定向
	if result.Error != nil {
		c.String(404, "未找到短链接")
		return
	}
	// 更新最后访问时间
	shortLink.LastRequestTime = time.Now().Unix()
	database.DB.Save(&shortLink)
	uri := config.Default.BasePath + shortLink.Url
	c.Redirect(http.StatusTemporaryRedirect, uri)
}
