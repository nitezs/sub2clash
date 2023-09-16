package utils

import (
	"sort"
	"strings"
	"sub2clash/model"
	"sub2clash/parser"
)

func GetContryName(proxy model.Proxy) string {
	// 创建一个切片包含所有的国家映射
	countryMaps := []map[string]string{
		model.CountryFlag,
		model.CountryChineseName,
		model.CountryISO,
		model.CountryEnglishName,
	}

	// 对每一个映射进行检查
	for _, countryMap := range countryMaps {
		for k, v := range countryMap {
			if strings.Contains(proxy.Name, k) {
				return v
			}
		}
	}

	return "其他地区"
}

func AddProxy(
	sub *model.Subscription, autotest bool,
	lazy bool, sortStrategy string,
	clashType model.ClashType, proxies ...model.Proxy,
) {
	newCountryGroupNames := make([]string, 0)
	proxyTypes := model.GetSupportProxyTypes(clashType)

	// 添加节点
	for _, proxy := range proxies {
		if !proxyTypes[proxy.Type] {
			continue
		}
		sub.Proxies = append(sub.Proxies, proxy)
		haveProxyGroup := false
		countryName := GetContryName(proxy)
		for i := range sub.ProxyGroups {
			group := &sub.ProxyGroups[i]

			if group.Name == countryName {
				group.Proxies = append(group.Proxies, proxy.Name)
				group.Size++
				haveProxyGroup = true
			}

			if group.Name == "手动切换" {
				group.Proxies = append(group.Proxies, proxy.Name)
				group.Size++
			}
		}
		if !haveProxyGroup {
			var newGroup model.ProxyGroup
			if !autotest {
				newGroup = model.ProxyGroup{
					Name:          countryName,
					Type:          "select",
					Proxies:       []string{proxy.Name},
					IsCountryGrop: true,
					Size:          1,
				}
			} else {
				newGroup = model.ProxyGroup{
					Name:          countryName,
					Type:          "url-test",
					Proxies:       []string{proxy.Name},
					IsCountryGrop: true,
					Url:           "http://www.gstatic.com/generate_204",
					Interval:      300,
					Tolerance:     50,
					Lazy:          lazy,
					Size:          1,
				}
			}
			sub.ProxyGroups = append(sub.ProxyGroups, newGroup)
			newCountryGroupNames = append(newCountryGroupNames, countryName)
		}
	}
	// 统计国家策略组数量
	countryGroupCount := 0
	for i := range sub.ProxyGroups {
		if sub.ProxyGroups[i].IsCountryGrop {
			countryGroupCount++
		}
	}
	// 对国家策略组进行排序
	switch sortStrategy {
	case "sizeasc":
		sort.Sort(model.ProxyGroupsSortBySize(sub.ProxyGroups[:countryGroupCount]))
	case "sizedesc":
		sort.Sort(sort.Reverse(model.ProxyGroupsSortBySize(sub.ProxyGroups[:countryGroupCount])))
	case "nameasc":
		sort.Sort(model.ProxyGroupsSortByName(sub.ProxyGroups[:countryGroupCount]))
	case "namedesc":
		sort.Sort(sort.Reverse(model.ProxyGroupsSortByName(sub.ProxyGroups[:countryGroupCount])))
	default:
		sort.Sort(model.ProxyGroupsSortByName(sub.ProxyGroups[:countryGroupCount]))
	}
}

func ParseProxy(proxies ...string) []model.Proxy {
	var result []model.Proxy
	for _, proxy := range proxies {
		if proxy != "" {
			var proxyItem model.Proxy
			var err error
			// 解析节点
			if strings.HasPrefix(proxy, "ss://") {
				proxyItem, err = parser.ParseSS(proxy)
			}
			if strings.HasPrefix(proxy, "trojan://") {
				proxyItem, err = parser.ParseTrojan(proxy)
			}
			if strings.HasPrefix(proxy, "vmess://") {
				proxyItem, err = parser.ParseVmess(proxy)
			}
			if strings.HasPrefix(proxy, "vless://") {
				proxyItem, err = parser.ParseVless(proxy)
			}
			if strings.HasPrefix(proxy, "ssr://") {
				proxyItem, err = parser.ParseShadowsocksR(proxy)
			}
			if err == nil {
				result = append(result, proxyItem)
			}
		}
	}
	return result
}
