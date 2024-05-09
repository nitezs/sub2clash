package handler

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"sub2clash/common"
	"sub2clash/common/database"
	"sub2clash/config"
	"sub2clash/model"
	"sub2clash/validator"

	"github.com/gin-gonic/gin"
)

func respondWithError(c *gin.Context, code int, message string) {
	c.String(code, message)
	c.Abort()
}

func GenerateLinkHandler(c *gin.Context) {
	var params validator.ShortLinkGenValidator
	if err := c.ShouldBind(&params); err != nil {
		respondWithError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if strings.TrimSpace(params.Url) == "" {
		respondWithError(c, http.StatusBadRequest, "URL 不能为空")
		return
	}

	hash, err := generateUniqueHash()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "生成短链接失败")
		return
	}

	shortLink := model.ShortLink{
		Hash:     hash,
		Url:      params.Url,
		Password: params.Password,
	}

	if err := database.SaveShortLink(&shortLink); err != nil {
		respondWithError(c, http.StatusInternalServerError, "数据库错误")
		return
	}

	if params.Password != "" {
		hash += "?password=" + params.Password
	}
	c.String(http.StatusOK, hash)
}

func generateUniqueHash() (string, error) {
	for {
		hash := common.RandomString(config.Default.ShortLinkLength)
		exists, err := database.CheckShortLinkHashExists(hash)
		if err != nil {
			return "", err
		}
		if !exists {
			return hash, nil
		}
	}
}

func UpdateLinkHandler(c *gin.Context) {
	var params validator.ShortLinkUpdateValidator
	if err := c.ShouldBindJSON(&params); err != nil {
		respondWithError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	shortLink := model.ShortLink{
		Hash:     params.Hash,
		Url:      params.Url,
		Password: params.Password,
	}
	if err := database.SaveShortLink(&shortLink); err != nil {
		respondWithError(c, http.StatusInternalServerError, "数据库错误")
		return
	}

	c.String(http.StatusOK, "短链接更新成功")
}

func GetRawConfHandler(c *gin.Context) {
	// 获取动态路由参数
	hash := c.Param("hash")
	password := c.Query("password")

	if strings.TrimSpace(hash) == "" {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	// 查询数据库中的短链接
	shortLink, err := database.FindShortLinkByHash(hash)
	if err != nil {
		c.String(http.StatusNotFound, "未找到短链接或密码错误")
		return
	}

	// 校验密码
	if shortLink.Password != "" && shortLink.Password != password {
		c.String(http.StatusNotFound, "未找到短链接或密码错误")
		return
	}

	// 更新最后访问时间
	shortLink.LastRequestTime = time.Now().Unix()
	err = database.SaveShortLink(shortLink)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "数据库错误")
		return
	}
	// 请求短链接指向的URL
	response, err := http.Get("http://localhost:" + strconv.Itoa(config.Default.Port) + "/" + shortLink.Url)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "请求错误: "+err.Error())
		return
	}
	defer response.Body.Close()

	// 读取响应内容
	all, err := io.ReadAll(response.Body)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "读取错误: "+err.Error())
		return
	}

	// 返回响应内容
	c.String(http.StatusOK, string(all))
}

func GetRawConfUriHandler(c *gin.Context) {
	// 获取动态路由参数
	hash := c.Query("hash")
	password := c.Query("password")

	if strings.TrimSpace(hash) == "" {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	// 查询数据库中的短链接
	shortLink, err := database.FindShortLinkByHash(hash)
	if err != nil {
		c.String(http.StatusNotFound, "未找到短链接或密码错误")
		return
	}

	// 校验密码
	if shortLink.Password != "" && shortLink.Password != password {
		c.String(http.StatusNotFound, "未找到短链接或密码错误")
		return
	}

	c.String(http.StatusOK, shortLink.Url)
}
