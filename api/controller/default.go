package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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
	groupsParam := parsedURL.Query().Get("subTags")

	// 分割字符串并创建map, 并插入每个proxy结构中
	if groupsParam != "" {
		subTags := strings.Split(groupsParam, ",")
		for i := range newProxies {
			newProxies[i].SubTags = []string{}
			for _, group := range subTags {
				newProxies[i].SubTags = append(newProxies[i].SubTags, group)
			}
		}
	}
}

func ParseCountries(newProxies []model.Proxy) {
	for i := range newProxies {
		countryName := utils.GetContryName(newProxies[i].Name)
		newProxies[i].Country = countryName
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
		//	SubTags           	[]string
		//  ...
		// 参考代码：
		//	if subName != "" {
		//	for i := range newProxies {
		//		newProxies[i].SubName = subName
		//	}
		//}
		ParseGroupTags(query.Subs[i], newProxies)

		ParseCountries(newProxies)

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
	MergeSubAndTemplate(temp, t, query.AutoTest, query.Lazy)
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

type Condition interface{}

type LogicCondition struct {
	And []Condition `json:"$and,omitempty"`
	Or  []Condition `json:"$or,omitempty"`
	Not []Condition `json:"$not,omitempty"`
}

type ComparisonCondition struct {
	Eq    string `json:"$eq,omitempty"`
	Regex string `json:"$regex,omitempty"`
}

type FieldCondition map[string]ComparisonCondition

func parseQuery(syntax string) (Condition, error) {
	if syntax == "{}" {
		return "{}", nil
	}
	var condition Condition
	err := json.Unmarshal([]byte(syntax), &condition)
	if err != nil {
		return nil, fmt.Errorf("error parsing query: %v", err)
	}
	return condition, nil
}

func matchProxy(proxy model.Proxy, condition Condition) bool {
	// 检查是否为匹配所有节点的特殊条件
	if cond, ok := condition.(string); ok && cond == "{}" {
		return true
	}
	switch cond := condition.(type) {
	case map[string]interface{}:
		return matchMapCondition(proxy, cond)
	default:
		return false
	}
}

// 新增一个辅助函数，用于从代理模型中提取数组类型的字段
func getProxyFieldArray(proxy model.Proxy, field string) []string {
	// 根据field的名称，从proxy对象中提取数组类型的字段值
	// 示例: 这里假设model.Proxy有某个字段是字符串数组
	switch field {
	case "SubTags":
		return proxy.SubTags
	// ... 其他可能的数组字段
	default:
		return nil
	}
}

func matchMapCondition(proxy model.Proxy, condition map[string]interface{}) bool {
	// Handle logic conditions ($and, $or, $not)
	if andConditions, ok := condition["$and"].([]interface{}); ok {
		for _, andCond := range andConditions {
			if !matchProxy(proxy, andCond) {
				return false
			}
		}
		return true
	}

	if orConditions, ok := condition["$or"].([]interface{}); ok {
		for _, orCond := range orConditions {
			if matchProxy(proxy, orCond) {
				return true
			}
		}
		return false
	}

	if notConditions, ok := condition["$not"].([]interface{}); ok {
		for _, notCond := range notConditions {
			if matchProxy(proxy, notCond) {
				return false
			}
		}
		return true
	}

	// Handle field conditions ($eq, $regex)
	for field, compCond := range condition {
		if compCondMap, ok := compCond.(map[string]interface{}); ok {
			if eqValue, ok := compCondMap["$eq"]; ok {
				if getProxyFieldValue(proxy, field) != eqValue {
					return false
				}
			}
			if regexValue, ok := compCondMap["$regex"]; ok {
				matched, _ := regexp.MatchString(regexValue.(string), getProxyFieldValue(proxy, field))
				if !matched {
					return false
				}
			}
		}
	}

	for field, compCond := range condition {
		if elemMatchMap, ok := compCond.(map[string]interface{}); ok {
			if elemMatchCond, ok := elemMatchMap["$elemMatch"]; ok {
				// 提取数组字段
				fieldValues := getProxyFieldArray(proxy, field)
				// 检查数组中是否有任何元素符合条件
				for _, fieldValue := range fieldValues {
					if matchProxyFieldValue(fieldValue, elemMatchCond.(map[string]interface{})) {
						return true
					}
				}
				return false
			}
		}
	}

	return true
}

// 专门用于检查单个字段值是否符合条件的函数
func matchProxyFieldValue(fieldValue string, condition map[string]interface{}) bool {
	for key, value := range condition {
		switch key {
		case "$eq":
			if fieldValue != value {
				return false
			}
		case "$regex":
			matched, _ := regexp.MatchString(value.(string), fieldValue)
			if !matched {
				return false
			}
		}
	}
	return true
}

func getProxyFieldValue(proxy model.Proxy, field string) string {
	// 根据field的名称，从proxy对象中提取相应的值
	// 示例: 这里假设model.Proxy有Name和Country等字段
	switch field {
	case "Name":
		return proxy.Name
	case "Country":
		return proxy.Country
	// ... 其他字段
	default:
		return ""
	}
}

func parseSyntaxA(syntax string, sub *model.Subscription) ([]string, []model.Proxy) {
	// 这里应该实现语法A的解析逻辑
	// 根据语法A的规则，过滤并返回匹配的代理名称列表
	// 示例代码仅作为逻辑框架，并非完整实现
	// 例如，你可能需要实现函数parseQuery, matchProxy等来处理复杂的查询逻辑

	// 假设解析后的查询结构体
	query, _ := parseQuery(syntax)

	// 过滤并返回匹配的代理名称
	matchedProxyNames := make([]string, 0)
	matchedProxies := make([]model.Proxy, 0)
	for _, proxy := range sub.Proxies {
		if matchProxy(proxy, query) {
			matchedProxyNames = append(matchedProxyNames, proxy.Name)
			matchedProxies = append(matchedProxies, proxy)
		}
	}
	return matchedProxyNames, matchedProxies
}

// 添加到某个节点组
func AddNewGroup(sub *model.Subscription, insertGroup string, autotest bool, lazy bool) {
	var newGroup model.ProxyGroup
	if !autotest {
		newGroup = model.ProxyGroup{
			Name:          insertGroup,
			Type:          "select",
			Proxies:       []string{},
			IsCountryGrop: true,
			Size:          1,
		}
	} else {
		newGroup = model.ProxyGroup{
			Name:          insertGroup,
			Type:          "url-test",
			Proxies:       []string{},
			IsCountryGrop: true,
			Url:           "https://www.gstatic.com/generate_204",
			Interval:      300,
			Tolerance:     50,
			Lazy:          lazy,
			Size:          1,
		}
	}
	sub.ProxyGroups = append(sub.ProxyGroups, newGroup)
}

// 添加到某个节点组
func AddToGroup(sub *model.Subscription, proxy model.Proxy, insertGroup string) bool {
	for i := range sub.ProxyGroups {
		group := &sub.ProxyGroups[i]

		if group.Name == insertGroup {
			group.Proxies = append(group.Proxies, proxy.Name)
			group.Size++
			return true
		}
	}
	return false
}

func MergeSubAndTemplate(temp *model.Subscription, sub *model.Subscription, autotest bool, lazy bool) {
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
			proxyName := temp.ProxyGroups[i].Proxies[j]
			if strings.HasPrefix(proxyName, "<") && strings.HasSuffix(proxyName, ">") {
				// 解析语法A
				syntax := strings.Trim(proxyName, "<>")
				proxyNames, proxies := parseSyntaxA(syntax, sub)

				// 把proxies放到一个新组中 L
				for index, _ := range proxyNames {
					// 遍历节点组，看是否有当前国家的组，如果没有，则新增，同时将
					var insertSuccess = AddToGroup(sub, proxies[index], proxyName)

					// 如果不存在此节点组，需要新增
					if !insertSuccess {
						AddNewGroup(sub, proxyName, autotest, lazy)
						// 同时将新节点插入到组中
						var _ = AddToGroup(sub, proxies[index], proxyName)
					}
				}

				newProxies = append(newProxies, proxyName)
			} else {
				newProxies = append(newProxies, proxyName)
			}
		}
		temp.ProxyGroups[i].Proxies = newProxies
	}
	temp.ProxyGroups = append(temp.ProxyGroups, sub.ProxyGroups...)
}
