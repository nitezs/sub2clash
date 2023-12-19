package parser

import (
	"net/url"
	"strings"
	"sub2clash/model"
)

func ParsePlugin(proxyURL string) (*model.Proxy, error) {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	proxy := &model.Proxy{}
	queryParams, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}

	if plugin, ok := queryParams["plugin"]; ok && len(plugin) > 0 {
		pluginOpts := make(map[string]any)
		pluginParts := strings.Split(plugin[0], ";")
		if len(pluginParts) > 0 {
			// 第一个部分是插件名称，我们需要特别处理
			for i, part := range pluginParts {
				if i == 0 {
					if part == "obfs-local" {
						proxy.Plugin = "obfs"
					}
				} else {
					opt := strings.SplitN(part, "=", 2)
					if len(opt) == 2 {
						if opt[0] == "obfs" {
							pluginOpts["mode"] = opt[1]
						} else if opt[0] == "obfs-host" {
							pluginOpts["host"] = opt[1]
						} else if opt[0] == "tfo" && opt[1] == "1" {
							proxy.Tfo = true
						}
					}
				}
			}
		}
		proxy.PluginOpts = pluginOpts
	}

	return proxy, nil
}
