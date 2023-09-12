package controller

import (
	"errors"
	"gopkg.in/yaml.v3"
	"net/url"
	"strings"
	"sub2clash/model"
	"sub2clash/parser"
	"sub2clash/utils"
)

func MixinSubsAndTemplate(subs []string, refresh bool, template string) (
	*model.Subscription, error,
) {
	// 定义变量
	var temp *model.Subscription
	var sub *model.Subscription
	// 加载模板
	template, err := utils.LoadTemplate(template)
	if err != nil {
		return nil, errors.New("加载模板失败: " + err.Error())
	}
	// 解析模板
	err = yaml.Unmarshal([]byte(template), &temp)
	if err != nil {
		return nil, errors.New("解析模板失败: " + err.Error())
	}
	var proxies []model.Proxy
	// 加载订阅
	for i := range subs {
		subs[i], _ = url.QueryUnescape(subs[i])
		if _, err := url.ParseRequestURI(subs[i]); err != nil {
			return nil, errors.New("订阅地址错误: " + err.Error())
		}
		data, err := utils.LoadSubscription(
			subs[i],
			refresh,
		)
		if err != nil {
			return nil, errors.New("加载订阅失败: " + err.Error())
		}
		// 解析订阅
		var proxyList []model.Proxy
		err = yaml.Unmarshal(data, &sub)
		if err != nil {
			// 如果无法直接解析，尝试Base64解码
			base64, err := parser.DecodeBase64(string(data))
			if err != nil {
				return nil, errors.New("加载订阅失败: " + err.Error())
			}
			proxyList = utils.ParseProxy(strings.Split(base64, "\n")...)
		} else {
			proxyList = sub.Proxies
		}
		proxies = append(proxies, proxyList...)
	}
	// 添加节点
	utils.AddProxy(temp, proxies...)
	return temp, nil
}
