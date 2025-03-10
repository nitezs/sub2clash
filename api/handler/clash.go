package handler

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/nitezs/sub2clash/config"
	"github.com/nitezs/sub2clash/model"
	"github.com/nitezs/sub2clash/validator"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func SubmodHandler(c *gin.Context) {

	query, err := validator.ParseQuery(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// 默认使用clash请求头
	if query.UserAgent == "" {
		query.UserAgent = "clash.meta/mihomo"
	}

	sub, err := BuildSub(model.Clash, query, config.Default.ClashTemplate)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if query.NodeListMode {
		nodelist := model.NodeList{}
		nodelist.Proxies = sub.Proxies
		marshal, err := yaml.Marshal(nodelist)
		if err != nil {
			c.String(http.StatusInternalServerError, "YAML序列化失败: "+err.Error())
			return
		}
		c.String(http.StatusOK, string(marshal))
		return
	}
	marshal, err := yaml.Marshal(sub)
	if err != nil {
		c.String(http.StatusInternalServerError, "YAML序列化失败: "+err.Error())
		return
	}
	// 如果有订阅名则设置
	if userAgent := c.GetHeader("User-Agent"); sub.SubscriptionName != "" && strings.Contains(userAgent, "clash") {
		c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+url.QueryEscape(sub.SubscriptionName))
	}
	c.String(http.StatusOK, string(marshal))
}
