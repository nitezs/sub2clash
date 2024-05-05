package parser

import (
	"net/url"
	"strings"
	"sub2clash/constant"
	"sub2clash/model"
)

func ParseVless(proxy string) (model.Proxy, error) {
	if !strings.HasPrefix(proxy, constant.VLESSPrefix) {
		return model.Proxy{}, &ParseError{Type: ErrInvalidPrefix, Raw: proxy}
	}

	urlParts := strings.SplitN(strings.TrimPrefix(proxy, constant.VLESSPrefix), "@", 2)
	if len(urlParts) != 2 {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidStruct,
			Message: "missing character '@' in url",
			Raw:     proxy,
		}
	}

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
	port, err := ParsePort(portStr)
	if err != nil {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidPort,
			Message: err.Error(),
			Raw:     proxy,
		}
	}

	params, err := url.ParseQuery(serverAndPortAndParams[1])
	if err != nil {
		return model.Proxy{}, &ParseError{
			Type:    ErrCannotParseParams,
			Raw:     proxy,
			Message: err.Error(),
		}
	}

	remarks := ""
	if len(serverInfo) == 2 {
		if strings.Contains(serverInfo[1], "|") {
			remarks = strings.SplitN(serverInfo[1], "|", 2)[1]
		} else {
			remarks, err = url.QueryUnescape(serverInfo[1])
			if err != nil {
				return model.Proxy{}, &ParseError{
					Type:    ErrCannotParseParams,
					Raw:     proxy,
					Message: err.Error(),
				}
			}
		}
	} else {
		remarks, err = url.QueryUnescape(server)
		if err != nil {
			return model.Proxy{}, err
		}
	}

	uuid := strings.TrimSpace(urlParts[0])
	flow, security, alpnStr, sni, insecure, fp, pbk, sid, path, host, serviceName, _type := params.Get("flow"), params.Get("security"), params.Get("alpn"), params.Get("sni"), params.Get("allowInsecure"), params.Get("fp"), params.Get("pbk"), params.Get("sid"), params.Get("path"), params.Get("host"), params.Get("serviceName"), params.Get("type")

	// enableUTLS := fp != ""
	insecureBool := insecure == "1"
	var alpn []string
	if strings.Contains(alpnStr, ",") {
		alpn = strings.Split(alpnStr, ",")
	} else {
		alpn = nil
	}

	result := model.Proxy{
		Type:   "vless",
		Server: server,
		Name:   remarks,
		Port:   port,
		UUID:   uuid,
		Flow:   flow,
	}

	if security == "tls" {
		result.TLS = true
		result.Alpn = alpn
		result.Sni = sni
		result.AllowInsecure = insecureBool
		result.ClientFingerprint = fp
	}

	if security == "reality" {
		result.TLS = true
		result.Servername = sni
		result.RealityOpts = model.RealityOptions{
			PublicKey: pbk,
			ShortID:   sid,
		}
		result.ClientFingerprint = fp
	}

	if _type == "ws" {
		result.Network = "ws"
		result.WSOpts = model.WSOptions{
			Path: path,
		}
		if host != "" {
			result.WSOpts.Headers = make(map[string]string)
			result.WSOpts.Headers["Host"] = host
		}
	}

	// if _type == "quic" {
	// 	// 未查到相关支持文档
	// }

	if _type == "grpc" {
		result.Network = "grpc"
		result.GrpcOpts = model.GrpcOptions{
			GrpcServiceName: serviceName,
		}
	}

	if _type == "http" {
		hosts, err := url.QueryUnescape(host)
		if err != nil {
			return model.Proxy{}, &ParseError{
				Type:    ErrCannotParseParams,
				Raw:     proxy,
				Message: err.Error(),
			}
		}
		result.Network = "http"
		result.HTTPOpts = model.HTTPOptions{
			Headers: map[string][]string{"Host": strings.Split(hosts, ",")},
		}

	}

	return result, nil
}
