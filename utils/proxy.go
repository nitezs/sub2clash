package utils

import (
	"strings"
	"sub2clash/logger"
	"sub2clash/model"
	"sub2clash/parser"

	"go.uber.org/zap"
)

func GetContryName(countryKey string) string {
	// 创建一个切片包含所有的国家映射
	countryMaps := []map[string]string{
		model.CountryFlag,
		model.CountryChineseName,
		model.CountryISO,
		model.CountryEnglishName,
	}

	checkForTW := false
	// 对每一个映射进行检查
	for i, countryMap := range countryMaps {
		if i == 2 {
			// 对ISO匹配做特殊处理
			// 根据常用分割字符分割字符串
			splitChars := []string{"-", "_", " "}
			key := make([]string, 0)
			for _, splitChar := range splitChars {
				slic := strings.Split(countryKey, splitChar)
				for _, v := range slic {
					if len(v) == 2 {
						key = append(key, v)
					}
				}
			}
			// 对每一个分割后的字符串进行检查
			for _, v := range key {
				// 如果匹配到了国家
				if country, ok := countryMap[strings.ToUpper(v)]; ok {
					v := country

					if v == "中国(CN)" {
						checkForTW = true
						continue
					}
					if checkForTW && v == "台湾(TW)" {
						// 如果正在检测是否是台湾
						return v
					}

					// 台湾可能误判成CN了，先不返回，等待之后确认不是台湾
					return v
				}
			}
		}
		for k, v := range countryMap {
			if strings.Contains(countryKey, k) {
				if v == "中国(CN)" {
					checkForTW = true
					continue
				}
				if checkForTW && v == "台湾(TW)" {
					// 如果正在检测是否是台湾
					return v
				}

				// 台湾可能误判成CN了，先不返回，等待之后确认不是台湾
				return v
			}
		}
	}
	if checkForTW {
		return "中国(CN)"
	}
	return "其他地区"
}

func AddAllNewProxies(
	sub *model.Subscription, lazy bool, clashType model.ClashType, proxies ...model.Proxy,
) {
	proxyTypes := model.GetSupportProxyTypes(clashType)

	// 遍历每个代理节点，添加节点
	for _, proxy := range proxies {
		// 跳过无效类型
		if !proxyTypes[proxy.Type] {
			continue
		}
		sub.Proxies = append(sub.Proxies, proxy)
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
			if strings.HasPrefix(proxy, "hysteria2://") {
				proxyItem, err = parser.ParseHysteria2(proxy)
			}
			if err == nil {
				// todo: 解析plugin字段，包括plugin和opt 填充到model.Proxy结构中，并与proxyItem合并
				pluginProxyItem, err := parser.ParsePlugin(proxy)
				if err == nil {
					proxyItem.Plugin = pluginProxyItem.Plugin
					proxyItem.PluginOpts = pluginProxyItem.PluginOpts
				}
				result = append(result, proxyItem)
			} else {
				logger.Logger.Debug(
					"parse proxy failed", zap.String("proxy", proxy), zap.Error(err),
				)
			}
		}
	}
	return result
}
