package utils

import (
	"sort"
	"strings"
	"sub/model"
	"sub/parser"
)

func GetContryCode(proxy model.Proxy) string {
	keys := make([]string, 0, len(model.CountryKeywords))
	for k := range model.CountryKeywords {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if strings.Contains(strings.ToLower(proxy.Name), strings.ToLower(k)) {
			return model.CountryKeywords[k]
		}
	}
	return "其他地区"
}

var skipGroups = map[string]bool{
	"手动切换": true,
	"全球直连": true,
	"广告拦截": true,
	"应用进化": true,
}

func AddProxy(sub *model.Subscription, proxies ...model.Proxy) {
	newContryNames := make([]string, 0, len(proxies))
	for p := range proxies {
		proxy := proxies[p]
		sub.Proxies = append(sub.Proxies, proxy)
		haveProxyGroup := false
		for i := range sub.ProxyGroups {
			group := &sub.ProxyGroups[i]
			groupName := []rune(group.Name)
			proxyName := []rune(proxy.Name)

			if string(groupName[:2]) == string(proxyName[:2]) || GetContryCode(proxy) == group.Name {
				group.Proxies = append(group.Proxies, proxy.Name)
				haveProxyGroup = true
			}

			if group.Name == "手动切换" {
				group.Proxies = append(group.Proxies, proxy.Name)
			}
		}
		if !haveProxyGroup {
			contryCode := GetContryCode(proxy)
			newGroup := model.ProxyGroup{
				Name:    contryCode,
				Type:    "select",
				Proxies: []string{proxy.Name},
			}
			newContryNames = append(newContryNames, contryCode)
			sub.ProxyGroups = append(sub.ProxyGroups, newGroup)
		}
	}
	newContryNamesMap := make(map[string]bool)
	for _, n := range newContryNames {
		newContryNamesMap[n] = true
	}
	for i := range sub.ProxyGroups {
		if !skipGroups[sub.ProxyGroups[i].Name] && !newContryNamesMap[sub.ProxyGroups[i].Name] {
			newProxies := make(
				[]string, len(newContryNames), len(newContryNames)+len(sub.ProxyGroups[i].Proxies),
			)
			copy(newProxies, newContryNames)
			sub.ProxyGroups[i].Proxies = append(newProxies, sub.ProxyGroups[i].Proxies...)
		}
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
