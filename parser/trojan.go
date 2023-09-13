package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sub2clash/model"
)

func ParseTrojan(proxy string) (model.Proxy, error) {
	// 判断是否以 trojan:// 开头
	if !strings.HasPrefix(proxy, "trojan://") {
		return model.Proxy{}, fmt.Errorf("无效的 trojan Url")
	}
	// 分割
	parts := strings.SplitN(strings.TrimPrefix(proxy, "trojan://"), "@", 2)
	if len(parts) != 2 {
		return model.Proxy{}, fmt.Errorf("无效的 trojan Url")
	}
	// 分割
	serverInfo := strings.SplitN(parts[1], "#", 2)
	serverAndPortAndParams := strings.SplitN(serverInfo[0], "?", 2)
	serverAndPort := strings.SplitN(serverAndPortAndParams[0], ":", 2)
	params, err := url.ParseQuery(serverAndPortAndParams[1])
	if err != nil {
		return model.Proxy{}, err
	}
	if len(serverAndPort) != 2 {
		return model.Proxy{}, fmt.Errorf("无效的 trojan 服务器和端口")
	}
	// 处理端口
	port, err := strconv.Atoi(strings.TrimSpace(serverAndPort[1]))
	if err != nil {
		return model.Proxy{}, err
	}
	// 返回结果
	result := model.Proxy{
		Type:     "trojan",
		Server:   strings.TrimSpace(serverAndPort[0]),
		Port:     port,
		UDP:      true,
		Password: strings.TrimSpace(parts[0]),
		Sni:      params.Get("sni"),
	}
	// 如果有节点名称
	if len(serverInfo) == 2 {
		result.Name, _ = url.QueryUnescape(strings.TrimSpace(serverInfo[1]))
	} else {
		result.Name = serverAndPort[0]
	}
	return result, nil
}
