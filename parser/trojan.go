package parser

import (
	"net/url"
	"strings"
	"sub2clash/constant"
	"sub2clash/model"
)

// ParseTrojan 解析给定的Trojan代理URL并返回Proxy结构。
func ParseTrojan(proxy string) (model.Proxy, error) {
	if !strings.HasPrefix(proxy, constant.TrojanPrefix) {
		return model.Proxy{}, &ParseError{Type: ErrInvalidPrefix, Raw: proxy}
	}

	proxy = strings.TrimPrefix(proxy, constant.TrojanPrefix)
	urlParts := strings.SplitN(proxy, "@", 2)
	if len(urlParts) != 2 {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidStruct,
			Message: "missing character '@' in url",
			Raw:     proxy,
		}
	}
	password := strings.TrimSpace(urlParts[0])

	serverInfo := strings.SplitN(urlParts[1], "#", 2)
	serverAndPortAndParams := strings.SplitN(serverInfo[0], "?", 2)
	if len(serverAndPortAndParams) != 2 {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidStruct,
			Message: "missing character '?' in url",
			Raw:     proxy,
		}
	}

	serverAndPort := strings.SplitN(serverAndPortAndParams[0], ":", 2)
	if len(serverAndPort) != 2 {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidStruct,
			Message: "missing server host or port",
			Raw:     proxy,
		}
	}
	server, portStr := serverAndPort[0], serverAndPort[1]

	params, err := url.ParseQuery(serverAndPortAndParams[1])
	if err != nil {
		return model.Proxy{}, &ParseError{
			Type:    ErrCannotParseParams,
			Raw:     proxy,
			Message: err.Error(),
		}
	}

	port, err := ParsePort(portStr)
	if err != nil {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidPort,
			Message: err.Error(),
			Raw:     proxy,
		}
	}

	remarks := ""
	if len(serverInfo) == 2 {
		remarks, _ = url.QueryUnescape(strings.TrimSpace(serverInfo[1]))
	} else {
		remarks = serverAndPort[0]
	}

	network, security, alpnStr, sni, pbk, sid, fp, path, host, serviceName := params.Get("type"), params.Get("security"), params.Get("alpn"), params.Get("sni"), params.Get("pbk"), params.Get("sid"), params.Get("fp"), params.Get("path"), params.Get("host"), params.Get("serviceName")

	var alpn []string
	if strings.Contains(alpnStr, ",") {
		alpn = strings.Split(alpnStr, ",")
	} else {
		alpn = nil
	}

	// enableUTLS := fp != ""

	// 构建Proxy结构体
	result := model.Proxy{
		Type:     "trojan",
		Server:   server,
		Port:     port,
		Password: password,
		Name:     remarks,
		Network:  network,
	}

	if security == "xtls" || security == "tls" {
		result.Alpn = alpn
		result.Sni = sni
		result.TLS = true
	}

	if security == "reality" {
		result.TLS = true
		result.Sni = sni
		result.RealityOpts = model.RealityOptions{
			PublicKey: pbk,
			ShortID:   sid,
		}
		result.Fingerprint = fp
	}

	if network == "ws" {
		result.Network = "ws"
		result.WSOpts = model.WSOptions{
			Path: path,
			Headers: map[string]string{
				"Host": host,
			},
		}
	}

	// if network == "http" {
	// 	// 未查到相关支持文档
	// }

	// if network == "quic" {
	// 	// 未查到相关支持文档
	// }

	if network == "grpc" {
		result.GrpcOpts = model.GrpcOptions{
			GrpcServiceName: serviceName,
		}
	}

	return result, nil
}
