package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sub2clash/model"
)

func ParseShadowsocksR(proxy string) (model.Proxy, error) {
	// 判断是否以 ssr:// 开头
	if !strings.HasPrefix(proxy, "ssr://") {
		return model.Proxy{}, fmt.Errorf("无效的 ssr Url")
	}
	var err error
	if !strings.Contains(proxy, ":") {
		proxy, err = DecodeBase64(strings.TrimPrefix(proxy, "ssr://"))
		if err != nil {
			return model.Proxy{}, err
		}
	}
	// 分割
	detailsAndParams := strings.SplitN(strings.TrimPrefix(proxy, "ssr://"), "/?", 2)
	parts := strings.Split(detailsAndParams[0], ":")
	params, err := url.ParseQuery(detailsAndParams[1])
	if err != nil {
		return model.Proxy{}, err
	}
	// 处理端口
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return model.Proxy{}, err
	}
	result := model.Proxy{
		Type:          "ssr",
		Server:        parts[0],
		Port:          port,
		Protocol:      parts[2],
		Cipher:        parts[3],
		Obfs:          parts[4],
		Password:      parts[5],
		ObfsParam:     params.Get("obfsparam"),
		ProtocolParam: params.Get("protoparam"),
	}
	return result, nil
}
