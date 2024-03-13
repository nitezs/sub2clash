package handler

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sub2clash/config"
	"sub2clash/logger"
	"sub2clash/model"
	"sub2clash/utils"
	"sub2clash/utils/database"
	"sub2clash/validator"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	var item model.ShortLink
	result := database.FindShortLinkByUrl(params.Url, &item)
	if result.Error == nil {
		if item.Password != params.Password {
			item.Password = params.Password
			database.SaveShortLink(&item)
			c.String(200, item.Hash+"?password="+params.Password)
		} else {
			c.String(200, item.Hash)
		}
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
		hash += "?password=" + params.Password
	}
	c.String(200, hash)
}

func ShortLinkGetUrlHandler(c *gin.Context) {
	var params validator.ShortLinkGetValidator
	if err := c.ShouldBindQuery(&params); err != nil {
		c.String(400, "参数错误: "+err.Error())
		return
	}
	if strings.TrimSpace(params.Hash) == "" {
		c.String(400, "参数错误")
		return
	}
	var shortLink model.ShortLink
	result := database.FindShortLinkByHash(params.Hash, &shortLink)
	if result.Error != nil {
		c.String(404, "未找到短链接")
		return
	}
	if shortLink.Password != "" && shortLink.Password != params.Password {
		c.String(403, "密码错误")
		return
	}
	c.String(200, shortLink.Url)
}

func ShortLinkGetConfigHandler(c *gin.Context) {
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
		c.String(404, "未找到短链接或密码错误")
		return
	}
	if shortLink.Password != "" && shortLink.Password != password {
		c.String(404, "未找到短链接或密码错误")
		return
	}
	// 更新最后访问时间
	shortLink.LastRequestTime = time.Now().Unix()
	database.SaveShortLink(&shortLink)
	get, err := utils.Get("http://localhost:" + strconv.Itoa(config.Default.Port) + "/" + shortLink.Url)
	if err != nil {
		logger.Logger.Debug("get short link data failed", zap.Error(err))
		c.String(500, "请求错误: "+err.Error())
		return
	}
	all, err := io.ReadAll(get.Body)
	if err != nil {
		logger.Logger.Debug("read short link data failed", zap.Error(err))
		c.String(500, "读取错误: "+err.Error())
		return
	}
	c.String(http.StatusOK, string(all))
}

func ShortLinkUpdateHandler(c *gin.Context) {
	var params validator.ShortLinkUpdateValidator
	if err := c.ShouldBind(&params); err != nil {
		c.String(400, "参数错误: "+err.Error())
	}
	if strings.TrimSpace(params.Url) == "" {
		c.String(400, "参数错误")
		return
	}
	var shortLink model.ShortLink
	result := database.FindShortLinkByHash(params.Hash, &shortLink)
	if result.Error != nil {
		c.String(404, "未找到短链接")
		return
	}
	if shortLink.Password == "" {
		c.String(403, "无法修改无密码短链接")
		return
	}
	if shortLink.Password != params.Password {
		c.String(403, "密码错误")
		return
	}
	shortLink.Url = params.Url
	database.SaveShortLink(&shortLink)
	c.String(200, "更新成功")
}
