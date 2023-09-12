package controller

import (
	"errors"
	"gopkg.in/yaml.v3"
	"strings"
	"sub/model"
	"sub/parser"
	"sub/utils"
	"sub/validator"
)

func MixinSubTemp(query validator.SubQuery, template string) (*model.Subscription, error) {
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
	// 加载订阅
	data, err := utils.LoadSubscription(
		query.Sub,
		query.Refresh,
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
	// 添加节点
	utils.AddProxy(temp, proxyList...)
	return temp, nil
}
