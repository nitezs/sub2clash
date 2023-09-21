package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sub2clash/logger"
	"sub2clash/model"
	"sub2clash/parser"
	"sub2clash/utils"
	"sub2clash/validator"
)

func BuildSub(clashType model.ClashType, query validator.SubValidator, template string) (
	*model.Subscription, error,
) {
	// 定义变量
	var temp = &model.Subscription{}
	var sub = &model.Subscription{}
	var err error
	var templateBytes []byte
	// 加载模板
	if query.Template != "" {
		template = query.Template
	}
	_, err = url.ParseRequestURI(template)
	if err != nil {
		templateBytes, err = utils.LoadTemplate(template)
		if err != nil {
			logger.Logger.Debug(
				"load template failed", zap.String("template", template), zap.Error(err),
			)
			return nil, errors.New("加载模板失败: " + err.Error())
		}
	} else {
		templateBytes, err = utils.LoadSubscription(template, query.Refresh)
		if err != nil {
			logger.Logger.Debug(
				"load template failed", zap.String("template", template), zap.Error(err),
			)
			return nil, errors.New("加载模板失败: " + err.Error())
		}
	}
	// 解析模板
	err = yaml.Unmarshal(templateBytes, &temp)
	if err != nil {
		logger.Logger.Debug("parse template failed", zap.Error(err))
		return nil, errors.New("解析模板失败: " + err.Error())
	}
	var proxyList []model.Proxy
	// 加载订阅
	for i := range query.Subs {
		data, err := utils.LoadSubscription(query.Subs[i], query.Refresh)
		if err != nil {
			logger.Logger.Debug(
				"load subscription failed", zap.String("url", query.Subs[i]), zap.Error(err),
			)
			return nil, errors.New("加载订阅失败: " + err.Error())
		}
		// 解析订阅
		err = yaml.Unmarshal(data, &sub)
		if err != nil {
			reg, _ := regexp.Compile("(ssr|ss|vmess|trojan|http|https)://")
			if reg.Match(data) {
				p := utils.ParseProxy(strings.Split(string(data), "\n")...)
				proxyList = append(proxyList, p...)
			} else {
				// 如果无法直接解析，尝试Base64解码
				base64, err := parser.DecodeBase64(string(data))
				if err != nil {
					logger.Logger.Debug(
						"parse subscription failed", zap.String("url", query.Subs[i]),
						zap.String("data", string(data)),
						zap.Error(err),
					)
					return nil, errors.New("加载订阅失败: " + err.Error())
				}
				p := utils.ParseProxy(strings.Split(base64, "\n")...)
				proxyList = append(proxyList, p...)
			}
		} else {
			proxyList = append(proxyList, sub.Proxies...)
		}
	}
	// 将新增节点都添加到临时变量 t 中，防止策略组排序错乱
	var t = &model.Subscription{}
	utils.AddProxy(t, query.AutoTest, query.Lazy, clashType, proxyList...)
	// 处理自定义代理
	utils.AddProxy(
		t, query.AutoTest, query.Lazy, clashType,
		utils.ParseProxy(query.Proxies...)...,
	)
	// 排序策略组
	switch query.Sort {
	case "sizeasc":
		sort.Sort(model.ProxyGroupsSortBySize(t.ProxyGroups))
	case "sizedesc":
		sort.Sort(sort.Reverse(model.ProxyGroupsSortBySize(t.ProxyGroups)))
	case "nameasc":
		sort.Sort(model.ProxyGroupsSortByName(t.ProxyGroups))
	case "namedesc":
		sort.Sort(sort.Reverse(model.ProxyGroupsSortByName(t.ProxyGroups)))
	default:
		sort.Sort(model.ProxyGroupsSortByName(t.ProxyGroups))
	}
	// 合并新节点和模板
	MergeSubAndTemplate(temp, t)
	// 处理自定义规则
	for _, v := range query.Rules {
		if v.Prepend {
			utils.PrependRules(temp, v.Rule)
		} else {
			utils.AppendRules(temp, v.Rule)
		}
	}
	// 处理自定义 ruleProvider
	for _, v := range query.RuleProviders {
		hash := sha256.Sum224([]byte(v.Url))
		name := hex.EncodeToString(hash[:])
		provider := model.RuleProvider{
			Type:     "http",
			Behavior: v.Behavior,
			Url:      v.Url,
			Path:     "./" + name + ".yaml",
			Interval: 3600,
		}
		if v.Prepend {
			utils.PrependRuleProvider(
				temp, v.Name, v.Group, provider,
			)
		} else {
			utils.AppenddRuleProvider(
				temp, v.Name, v.Group, provider,
			)
		}
	}
	return temp, nil
}

func MergeSubAndTemplate(temp *model.Subscription, sub *model.Subscription) {
	// 只合并节点、策略组
	// 统计所有国家策略组名称
	var countryGroupNames []string
	for _, proxyGroup := range sub.ProxyGroups {
		if proxyGroup.IsCountryGrop {
			countryGroupNames = append(
				countryGroupNames, proxyGroup.Name,
			)
		}
	}
	// 将订阅中的节点添加到模板中
	temp.Proxies = append(temp.Proxies, sub.Proxies...)
	// 将订阅中的策略组添加到模板中
	skipGroups := []string{"全球直连", "广告拦截", "手动切换"}
	for i := range temp.ProxyGroups {
		skip := false
		for _, v := range skipGroups {
			if strings.Contains(temp.ProxyGroups[i].Name, v) {
				if v == "手动切换" {
					proxies := make([]string, 0, len(sub.Proxies))
					for _, p := range sub.Proxies {
						proxies = append(proxies, p.Name)
					}
					temp.ProxyGroups[i].Proxies = proxies
				}
				skip = true
				continue
			}
		}
		if !skip {
			temp.ProxyGroups[i].Proxies = append(temp.ProxyGroups[i].Proxies, countryGroupNames...)
		}
	}
	temp.ProxyGroups = append(temp.ProxyGroups, sub.ProxyGroups...)
}
