package parser

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"sub2clash/model"
)

//hysteria://host:port?protocol=udp&auth=123456&peer=sni.domain&insecure=1&upmbps=100&downmbps=100&alpn=hysteria&obfs=xplus&obfsParam=123456#remarks
//
//- host: hostname or IP address of the server to connect to (required)
//- port: port of the server to connect to (required)
//- protocol: protocol to use ("udp", "wechat-video", "faketcp") (optional, default: "udp")
//- auth: authentication payload (string) (optional)
//- peer: SNI for TLS (optional)
//- insecure: ignore certificate errors (optional)
//- upmbps: upstream bandwidth in Mbps (required)
//- downmbps: downstream bandwidth in Mbps (required)
//- alpn: QUIC ALPN (optional)
//- obfs: Obfuscation mode (optional, empty or "xplus")
//- obfsParam: Obfuscation password (optional)
//- remarks: remarks (optional)

func ParseHysteria(proxy string) (model.Proxy, error) {
	// 判断是否以 hysteria:// 开头
	if !strings.HasPrefix(proxy, "hysteria://") {
		return model.Proxy{}, errors.New("invalid hysteria Url")
	}
	// 分割
	parts := strings.SplitN(strings.TrimPrefix(proxy, "hysteria://"), "?", 2)
	serverInfo := strings.SplitN(parts[0], ":", 2)
	if len(serverInfo) != 2 {
		return model.Proxy{}, errors.New("invalid hysteria Url")
	}
	params, err := url.ParseQuery(parts[1])
	if err != nil {
		return model.Proxy{}, errors.New("invalid hysteria Url")
	}
	host := serverInfo[0]
	port, err := strconv.Atoi(serverInfo[1])
	if err != nil {
		return model.Proxy{}, errors.New("invalid hysteria Url")
	}
	protocol := params.Get("protocol")
	auth := params.Get("auth")
	peer := params.Get("peer")
	insecure := params.Get("insecure")
	upmbps := params.Get("upmbps")
	downmbps := params.Get("downmbps")
	alpn := params.Get("alpn")
	obfs := params.Get("obfs")
	obfsParam := params.Get("obfsParam")
	remarks := ""
	if strings.Contains(parts[1], "#") {
		r := strings.Split(parts[1], "#")
		remarks = r[len(r)-1]
	} else {
		remarks = serverInfo[0] + ":" + serverInfo[1]
	}
	// 返回结果
	result := model.Proxy{
		Type:           "hysteria",
		Name:           remarks,
		Server:         host,
		Port:           port,
		Up:             upmbps,
		Down:           downmbps,
		Auth:           auth,
		Obfs:           obfs,
		Sni:            peer,
		SkipCertVerify: insecure == "1",
		Alpn:           strings.Split(alpn, ","),
		ObfsParam:      obfsParam,
		Protocol:       protocol,
	}
	return result, nil
}
