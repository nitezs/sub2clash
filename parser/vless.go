package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sub2clash/model"
)

func ParseVless(proxy string) (model.Proxy, error) {
	// 判断是否以 vless:// 开头
	if !strings.HasPrefix(proxy, "vless://") {
		return model.Proxy{}, fmt.Errorf("invalid vless Url")
	}
	// 分割
	parts := strings.SplitN(strings.TrimPrefix(proxy, "vless://"), "@", 2)
	if len(parts) != 2 {
		return model.Proxy{}, fmt.Errorf("invalid vless Url")
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
		return model.Proxy{}, fmt.Errorf("invalid vless")
	}
	// 处理端口
	port, err := strconv.Atoi(strings.TrimSpace(serverAndPort[1]))
	if err != nil {
		return model.Proxy{}, err
	}
	// 返回结果
	result := model.Proxy{
		Type:        "vless",
		Server:      strings.TrimSpace(serverAndPort[0]),
		Port:        port,
		UUID:        strings.TrimSpace(parts[0]),
		UDP:         true,
		Sni:         params.Get("sni"),
		Network:     params.Get("type"),
		TLS:         params.Get("security") == "tls",
		Flow:        params.Get("flow"),
		Fingerprint: params.Get("fp"),
		Servername:  params.Get("sni"),
		RealityOpts: model.RealityOptsStruct{
			PublicKey: params.Get("pbk"),
		},
	}
	if params.Get("alpn") != "" {
		result.Alpn = strings.Split(params.Get("alpn"), ",")
	}
	if params.Get("type") == "ws" {
		result.WSOpts = model.WSOptsStruct{
			Path: params.Get("path"),
			Headers: model.HeaderStruct{
				Host: params.Get("host"),
			},
		}
	}
	if params.Get("type") == "grpc" {
		result.GRPCOpts = model.GRPCOptsStruct{
			GRPCServiceName: params.Get("serviceName"),
		}
	}
	// 如果有节点名称
	if len(serverInfo) == 2 {
		if strings.Contains(serverInfo[1], "|") {
			result.Name = strings.SplitN(serverInfo[1], "|", 2)[1]
		} else {
			result.Name = serverInfo[1]
		}
	} else {
		result.Name = serverAndPort[0]
	}
	return result, nil
}
