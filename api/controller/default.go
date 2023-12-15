package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sub2clash/logger"
	"sub2clash/model"
	"sub2clash/parser"
	"sub2clash/utils"
	"sub2clash/validator"
)

func ParseGroupTags(subURL string, newProxies []model.Proxy) {
	parsedURL, _ := url.Parse(subURL)
	// 提取groups参数
	groupsParam := parsedURL.Query().Get("groups")

	// 分割字符串并创建map, 并插入每个proxy结构中
	if groupsParam != "" {
		groups := strings.Split(groupsParam, ",")
		for i := range newProxies {
			newProxies[i].GroupTags = make(map[string]bool)
			for _, group := range groups {
				newProxies[i].GroupTags[group] = true
			}
		}
	}
}

func WalkSubsForProxyList(sub *model.Subscription, query validator.SubValidator, proxyList *[]model.Proxy) (
	bool, error,
) {
	for i := range query.Subs {
		data, err := utils.LoadSubscription(query.Subs[i], query.Refresh)
		subName := ""
		if strings.Contains(query.Subs[i], "#") {
			subName = query.Subs[i][strings.LastIndex(query.Subs[i], "#")+1:]
		}
		if err != nil {
			logger.Logger.Debug(
				"load subscription failed", zap.String("url", query.Subs[i]), zap.Error(err),
			)
			return false, errors.New("加载订阅失败: " + err.Error())
		}
		// 解析订阅
		err = yaml.Unmarshal(data, &sub)
		newProxies := make([]model.Proxy, 0)
		if err != nil {
			reg, _ := regexp.Compile("(ssr|ss|vmess|trojan|vless)://")
			if reg.Match(data) {
				p := utils.ParseProxy(strings.Split(string(data), "\n")...)
				newProxies = p
			} else {
				// 如果无法直接解析，尝试Base64解码
				base64, err := parser.DecodeBase64(string(data))
				if err != nil {
					logger.Logger.Debug(
						"parse subscription failed", zap.String("url", query.Subs[i]),
						zap.String("data", string(data)),
						zap.Error(err),
					)
					return false, errors.New("加载订阅失败: " + err.Error())
				}
				p := utils.ParseProxy(strings.Split(base64, "\n")...)
				newProxies = p
			}
		} else {
			newProxies = sub.Proxies
		}
		if subName != "" {
			for i := range newProxies {
				newProxies[i].SubName = subName
			}
		}

		// 给每个节点添加订阅属性，从url中的groups寻找
		// 解析url变量query.Subs[i]，获取groups属性，groups属性使用,分割，建立一个go的map[string, true]结构放置这些groupName
		// newProxies:
		//  ...
		//	UDPOverTCPVersion   int            `yaml:"udp-over-tcp-version,omitempty"`
		//	SubName             string         `yaml:"-"`
		//	GroupTags           map[string]bool
		//  ...
		// 参考代码：
		//	if subName != "" {
		//	for i := range newProxies {
		//		newProxies[i].SubName = subName
		//	}
		//}
		ParseGroupTags(query.Subs[i], newProxies)

		*proxyList = append(*proxyList, newProxies...)
	}
	return true, nil
}

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
	_, err = url.ParseRequestURI(template) // 判断template是不是一个在线http的配置
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
	// 遍历订阅链接 获取 proxyList
	success, err := WalkSubsForProxyList(sub, query, &proxyList)
	if !success {
		return nil, err
	}
	// 添加自定义节点
	if len(query.Proxies) != 0 {
		proxyList = append(proxyList, utils.ParseProxy(query.Proxies...)...)
	}
	// 给节点添加订阅名称
	for i := range proxyList {
		if proxyList[i].SubName != "" {
			proxyList[i].Name = strings.TrimSpace(proxyList[i].SubName) + " " + strings.TrimSpace(proxyList[i].Name)
		}
	}
	// 去掉配置相同的节点
	proxies := make(map[string]*model.Proxy)
	newProxies := make([]model.Proxy, 0, len(proxyList))
	for i := range proxyList {
		key := proxyList[i].Server + ":" + strconv.Itoa(proxyList[i].Port) + ":" + proxyList[i].Type
		if _, exist := proxies[key]; !exist {
			proxies[key] = &proxyList[i]
			newProxies = append(newProxies, proxyList[i])
		}
	}
	proxyList = newProxies
	// 删除节点
	if strings.TrimSpace(query.Remove) != "" {
		newProxyList := make([]model.Proxy, 0, len(proxyList))
		for i := range proxyList {
			removeReg, err := regexp.Compile(query.Remove)
			if err != nil {
				logger.Logger.Debug("remove regexp compile failed", zap.Error(err))
				return nil, errors.New("remove 参数非法: " + err.Error())
			}
			// 删除匹配到的节点
			if removeReg.MatchString(proxyList[i].Name) {
				continue // 如果匹配到要删除的元素，跳过该元素，不添加到新切片中
			}
			newProxyList = append(newProxyList, proxyList[i]) // 将要保留的元素添加到新切片中
		}
		proxyList = newProxyList
	}
	// 重命名
	if len(query.ReplaceKeys) != 0 {
		// 创建重命名正则表达式
		replaceRegs := make([]*regexp.Regexp, 0, len(query.ReplaceKeys))
		for _, v := range query.ReplaceKeys {
			replaceReg, err := regexp.Compile(v)
			if err != nil {
				logger.Logger.Debug("replace regexp compile failed", zap.Error(err))
				return nil, errors.New("replace 参数非法: " + err.Error())
			}
			replaceRegs = append(replaceRegs, replaceReg)
		}
		for i := range proxyList {
			// 重命名匹配到的节点
			for j, v := range replaceRegs {
				if err != nil {
					logger.Logger.Debug("replace regexp compile failed", zap.Error(err))
					return nil, errors.New("replaceName 参数非法: " + err.Error())
				}
				if v.MatchString(proxyList[i].Name) {
					proxyList[i].Name = v.ReplaceAllString(
						proxyList[i].Name, query.ReplaceTo[j],
					)
				}
			}
		}
	}
	// 重名检测
	names := make(map[string]int)
	for i := range proxyList {
		if _, exist := names[proxyList[i].Name]; exist {
			proxyList[i].Name = proxyList[i].Name + " " + strconv.Itoa(names[proxyList[i].Name])
		}
		names[proxyList[i].Name] = names[proxyList[i].Name] + 1
	}
	// trim
	for i := range proxyList {
		proxyList[i].Name = strings.TrimSpace(proxyList[i].Name)
	}
	// 将新增节点都添加到临时变量 t 中，防止策略组排序错乱
	var t = &model.Subscription{}
	utils.AddAllNewProxies(t, query.AutoTest, query.Lazy, clashType, proxyList...)
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
	var proxyNames []string
	for _, proxy := range sub.Proxies {
		proxyNames = append(proxyNames, proxy.Name)
	}
	// 将订阅中的节点添加到模板中
	temp.Proxies = append(temp.Proxies, sub.Proxies...)
	// 将订阅中的策略组添加到模板中
	for i := range temp.ProxyGroups {
		if temp.ProxyGroups[i].IsCountryGrop {
			continue
		}
		newProxies := make([]string, 0)
		countryGroupMap := make(map[string]model.ProxyGroup)
		for _, v := range sub.ProxyGroups {
			if v.IsCountryGrop {
				countryGroupMap[v.Name] = v
			}
		}
		for j := range temp.ProxyGroups[i].Proxies {
			reg := regexp.MustCompile("<(.*?)>")
			if reg.Match([]byte(temp.ProxyGroups[i].Proxies[j])) {
				key := reg.FindStringSubmatch(temp.ProxyGroups[i].Proxies[j])[1]
				switch key {
				case "all":
					newProxies = append(newProxies, proxyNames...)
				case "countries":
					newProxies = append(newProxies, countryGroupNames...)
				default:
					if len(key) == 2 {
						newProxies = append(
							newProxies, countryGroupMap[utils.GetContryName(key)].Proxies...,
						)
					}
				}
			} else {
				newProxies = append(newProxies, temp.ProxyGroups[i].Proxies[j])
			}
		}
		temp.ProxyGroups[i].Proxies = newProxies
	}
	temp.ProxyGroups = append(temp.ProxyGroups, sub.ProxyGroups...)
}
