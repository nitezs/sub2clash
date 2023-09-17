package controller

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"gopkg.in/yaml.v3"
	"net/url"
	"regexp"
	"strings"
	"sub2clash/model"
	"sub2clash/parser"
	"sub2clash/utils"
	"sub2clash/validator"
)

func BuildSub(clashType model.ClashType, query validator.SubQuery, template string) (
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
			return nil, errors.New("加载模板失败: " + err.Error())
		}
	} else {
		templateBytes, err = utils.LoadSubscription(template, query.Refresh)
		if err != nil {
			return nil, errors.New("加载模板失败: " + err.Error())
		}
	}
	// 解析模板
	err = yaml.Unmarshal(templateBytes, &temp)
	if err != nil {
		return nil, errors.New("解析模板失败: " + err.Error())
	}
	// 加载订阅
	for i := range query.Subs {
		data, err := utils.LoadSubscription(query.Subs[i], query.Refresh)
		if err != nil {
			return nil, errors.New("加载订阅失败: " + err.Error())
		}
		// 解析订阅
		var proxyList []model.Proxy
		err = yaml.Unmarshal(data, &sub)
		if err != nil {
			reg, _ := regexp.Compile("(ssr|ss|vmess|trojan|http|https)://")
			if reg.Match(data) {
				proxyList = utils.ParseProxy(strings.Split(string(data), "\n")...)
			} else {
				// 如果无法直接解析，尝试Base64解码
				base64, err := parser.DecodeBase64(string(data))
				if err != nil {
					return nil, errors.New("加载订阅失败: " + err.Error())
				}
				proxyList = utils.ParseProxy(strings.Split(base64, "\n")...)
			}
		} else {
			proxyList = sub.Proxies
		}
		utils.AddProxy(sub, query.AutoTest, query.Lazy, query.Sort, clashType, proxyList...)
	}
	// 处理自定义代理
	utils.AddProxy(
		sub, query.AutoTest, query.Lazy, query.Sort, clashType,
		utils.ParseProxy(query.Proxies...)...,
	)
	MergeSubAndTemplate(temp, sub)
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
		hash := md5.Sum([]byte(v.Url))
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
