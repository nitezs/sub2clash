package controller

import (
	_ "embed"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"net/http"
	"sub2clash/config"
	"sub2clash/validator"
)

func SubHandler(c *gin.Context) {
	// 从请求中获取参数
	query, err := validator.ParseQuery(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	sub, err := BuildSub(query, config.Default.MetaTemplate)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	// 输出
	marshal, err := yaml.Marshal(sub)
	if err != nil {
		c.String(http.StatusInternalServerError, "YAML序列化失败: "+err.Error())
		return
	}
	c.String(http.StatusOK, string(marshal))
}
