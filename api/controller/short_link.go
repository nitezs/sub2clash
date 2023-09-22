package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
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
	if strings.TrimSpace(params.Url) == "" {
		c.String(400, "参数错误")
		return
	}
	// 生成hash
	hash := utils.RandomString(config.Default.ShortLinkLength)
	// 存入数据库
	var item model.ShortLink
	result := database.FindShortLinkByUrl(params.Url, &item)
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
	result = database.FindShortLinkByHash(hash, &item)
	for result.Error == nil {
		hash = utils.RandomString(config.Default.ShortLinkLength)
		result = database.FindShortLinkByHash(hash, &item)
	}
	// 创建记录
	database.FirstOrCreateShortLink(
		&model.ShortLink{
			Hash:            hash,
			Url:             params.Url,
			LastRequestTime: -1,
			Password:        params.Password,
		},
	)
	// 返回短链接
	if params.Password != "" {
		hash += "/?password=" + params.Password
	}
	c.String(200, hash)
}

func ShortLinkGetHandler(c *gin.Context) {
	// 获取动态路由
	hash := c.Param("hash")
	password := c.Query("password")
	if strings.TrimSpace(hash) == "" {
		c.String(400, "参数错误")
		return
	}
	// 查询数据库
	var shortLink model.ShortLink
	result := database.FindShortLinkByHash(hash, &shortLink)
	// 重定向
	if result.Error != nil {
		c.String(404, "未找到短链接")
		return
	}
	if shortLink.Password != "" && shortLink.Password != password {
		c.String(403, "密码错误")
		return
	}
	// 更新最后访问时间
	shortLink.LastRequestTime = time.Now().Unix()
	database.SaveShortLink(&shortLink)
	uri := config.Default.BasePath + shortLink.Url
	c.Redirect(http.StatusTemporaryRedirect, uri)
}
