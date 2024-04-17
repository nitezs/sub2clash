package parser

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"sub2clash/model"
)

// hysteria2://letmein@example.com/?insecure=1&obfs=salamander&obfs-password=gawrgura&pinSHA256=deadbeef&sni=real.example.com#name

func ParseHysteria2(proxy string) (model.Proxy, error) {
	// 判断是否以 hysteria2:// 开头
	if !strings.HasPrefix(proxy, "hysteria2://") && !strings.HasPrefix(proxy, "hy2://") {
		return model.Proxy{}, errors.New("invalid hysteria2 Url")
	}
	// 分割
	parts := strings.SplitN(strings.TrimPrefix(proxy, "hysteria2://"), "@", 2)
	// 分割
	serverInfo := strings.SplitN(parts[1], "/?", 2)
	serverAndPort := strings.SplitN(serverInfo[0], ":", 2)
	if len(serverAndPort) == 1 {
		serverAndPort = append(serverAndPort, "443")
	} else if len(serverAndPort) != 2 {
		return model.Proxy{}, errors.New("invalid hysteria2 Url")
	}
	params, err := url.ParseQuery(serverInfo[1])
	if err != nil {
		return model.Proxy{}, errors.New("invalid hysteria2 Url")
	}
	// 获取端口
	port, err := strconv.Atoi(serverAndPort[1])
	if err != nil {
		return model.Proxy{}, errors.New("invalid hysteria2 Url")
	}
	name := ""
	if strings.Contains(proxy, "#") {
		splitResult := strings.Split(proxy, "#")
		name, _ = url.QueryUnescape(splitResult[len(splitResult)-1])
	}
	// 返回结果
	result := model.Proxy{
		Type:           "hysteria2",
		Name:           name,
		Server:         serverAndPort[0],
		Port:           port,
		Password:       parts[0],
		Obfs:           params.Get("obfs"),
		ObfsParam:      params.Get("obfs-password"),
		Sni:            params.Get("sni"),
		SkipCertVerify: params.Get("insecure") == "1",
	}
	return result, nil
}
