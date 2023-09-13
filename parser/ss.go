package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sub2clash/model"
)

// ParseSS 解析 SS（Shadowsocks）Url
func ParseSS(proxy string) (model.Proxy, error) {
	// 判断是否以 ss:// 开头
	if !strings.HasPrefix(proxy, "ss://") {
		return model.Proxy{}, fmt.Errorf("无效的 ss Url")
	}
	// 分割
	parts := strings.SplitN(strings.TrimPrefix(proxy, "ss://"), "@", 2)
	if len(parts) != 2 {
		return model.Proxy{}, fmt.Errorf("无效的 ss Url")
	}
	if !strings.Contains(parts[0], ":") {
		// 解码
		decoded, err := DecodeBase64(parts[0])
		if err != nil {
			return model.Proxy{}, err
		}
		parts[0] = decoded
	}
	credentials := strings.SplitN(parts[0], ":", 2)
	if len(credentials) != 2 {
		return model.Proxy{}, fmt.Errorf("无效的 ss 凭证")
	}
	// 分割
	serverInfo := strings.SplitN(parts[1], "#", 2)
	serverAndPort := strings.SplitN(serverInfo[0], ":", 2)
	if len(serverAndPort) != 2 {
		return model.Proxy{}, fmt.Errorf("无效的 ss 服务器和端口")
	}
	// 转换端口字符串为数字
	port, err := strconv.Atoi(strings.TrimSpace(serverAndPort[1]))
	if err != nil {
		return model.Proxy{}, err
	}
	// 返回结果
	result := model.Proxy{
		Type:     "ss",
		Cipher:   strings.TrimSpace(credentials[0]),
		Password: strings.TrimSpace(credentials[1]),
		Server:   strings.TrimSpace(serverAndPort[0]),
		Port:     port,
		UDP:      true,
		Name:     serverAndPort[0],
	}
	// 如果有节点名称
	if len(serverInfo) == 2 {
		unescape, err := url.QueryUnescape(serverInfo[1])
		if err != nil {
			return model.Proxy{}, err
		}
		result.Name = strings.TrimSpace(unescape)
	} else {
		result.Name = strings.TrimSpace(serverAndPort[0])
	}

	return result, nil
}
