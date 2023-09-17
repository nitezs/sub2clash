package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
	"sub2clash/config"
	"sub2clash/model"
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
	// 生成短链接
	//hash := utils.RandomString(6)
	shortLink := sha256.Sum224([]byte(params.Url))
	hash := hex.EncodeToString(shortLink[:])
	// 存入数据库
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
	// 更新最后访问时间
	shortLink.LastRequestTime = time.Now().Unix()
	database.DB.Save(&shortLink)
	// 重定向
	if result.Error != nil {
		c.String(404, "未找到短链接")
		return
	}
	uri := config.Default.BasePath + shortLink.Url
	c.Redirect(http.StatusTemporaryRedirect, uri)
}
