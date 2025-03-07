package handler

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nitezs/sub2clash/common"
	"github.com/nitezs/sub2clash/common/database"
	"github.com/nitezs/sub2clash/config"
	"github.com/nitezs/sub2clash/model"
	"github.com/nitezs/sub2clash/validator"

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

	var hash string
	var err error
	
	if params.CustomID != "" {
		// 检查自定义ID是否已存在
		exists, err := database.CheckShortLinkHashExists(params.CustomID)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "数据库错误")
			return
		}
		if exists {
			respondWithError(c, http.StatusBadRequest, "自定义ID已存在")
			return
		}
		hash = params.CustomID
	} else {
		hash, err = generateUniqueHash()
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "生成短链接失败")
			return
		}
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

	// 先获取原有的短链接
	existingLink, err := database.FindShortLinkByHash(params.Hash)
	if err != nil {
		respondWithError(c, http.StatusNotFound, "未找到短链接")
		return
	}

	// 验证密码
	if existingLink.Password != params.Password {
		respondWithError(c, http.StatusUnauthorized, "密码错误")
		return
	}

	// 更新URL，但保持原密码不变
	shortLink := model.ShortLink{
		Hash:     params.Hash,
		Url:      params.Url,
		Password: existingLink.Password, // 保持原密码不变
	}

	if err := database.SaveShortLink(&shortLink); err != nil {
		respondWithError(c, http.StatusInternalServerError, "数据库错误")
		return
	}

	c.String(http.StatusOK, "短链接更新成功")
}

func GetRawConfHandler(c *gin.Context) {

	hash := c.Param("hash")
	password := c.Query("password")

	if strings.TrimSpace(hash) == "" {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	shortLink, err := database.FindShortLinkByHash(hash)
	if err != nil {
		c.String(http.StatusNotFound, "未找到短链接或密码错误")
		return
	}

	if shortLink.Password != "" && shortLink.Password != password {
		c.String(http.StatusNotFound, "未找到短链接或密码错误")
		return
	}

	shortLink.LastRequestTime = time.Now().Unix()
	err = database.SaveShortLink(shortLink)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "数据库错误")
		return
	}

	response, err := http.Get("http://localhost:" + strconv.Itoa(config.Default.Port) + "/" + shortLink.Url)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "请求错误: "+err.Error())
		return
	}
	defer response.Body.Close()

	all, err := io.ReadAll(response.Body)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "读取错误: "+err.Error())
		return
	}

	c.String(http.StatusOK, string(all))
}

func GetRawConfUriHandler(c *gin.Context) {

	hash := c.Query("hash")
	password := c.Query("password")

	if strings.TrimSpace(hash) == "" {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	shortLink, err := database.FindShortLinkByHash(hash)
	if err != nil {
		c.String(http.StatusNotFound, "未找到短链接或密码错误")
		return
	}

	if shortLink.Password != "" && shortLink.Password != password {
		c.String(http.StatusNotFound, "未找到短链接或密码错误")
		return
	}

	c.String(http.StatusOK, shortLink.Url)
}
