package controller

import (
	"net/http"
	"net/url"
	"sub/config"
	"sub/validator"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func SubmodHandler(c *gin.Context) {
	// 从请求中获取参数
	var query validator.SubQuery
	if err := c.ShouldBind(&query); err != nil {
		c.String(http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	query.Sub, _ = url.QueryUnescape(query.Sub)
	// 混合订阅和模板节点
	sub, err := MixinSubTemp(query, config.Default.ClashTemplate)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	// 添加自定义节点、规则
	// 输出
	bytes, err := yaml.Marshal(sub)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, string(bytes))
}
