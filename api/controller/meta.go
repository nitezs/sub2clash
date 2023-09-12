package controller

import (
	_ "embed"
	"net/http"
	"strings"
	"sub2clash/config"
	"sub2clash/validator"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func SubHandler(c *gin.Context) {
	// 从请求中获取参数
	var query validator.SubQuery
	if err := c.ShouldBind(&query); err != nil {
		c.String(http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	// 混合订阅和模板节点
	sub, err := MixinSubsAndTemplate(
		strings.Split(query.Sub, ","), query.Refresh, config.Default.MetaTemplate,
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	// 添加自定义节点、规则
	// 输出
	marshal, err := yaml.Marshal(sub)
	if err != nil {
		c.String(http.StatusInternalServerError, "YAML序列化失败: "+err.Error())
		return
	}
	c.String(http.StatusOK, string(marshal))
}
