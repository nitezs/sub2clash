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
		for _, part := range pluginParts {
			opt := strings.SplitN(part, "=", 2)
			if len(opt) == 2 {
				proxy.Plugin = opt[0]
				pluginOpts[opt[0]] = opt[1]
			}
		}
		proxy.PluginOpts = pluginOpts
	}

	return proxy, nil
}
