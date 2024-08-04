package parser

import (
	"net/url"
	"strings"
	"sub2clash/constant"
	"sub2clash/model"
)

func ParseHysteria2(proxy string) (model.Proxy, error) {
	if !strings.HasPrefix(proxy, constant.Hysteria2Prefix1) &&
		!strings.HasPrefix(proxy, constant.Hysteria2Prefix2) {
		return model.Proxy{}, &ParseError{Type: ErrInvalidPrefix, Raw: proxy}
	}

	proxy = strings.TrimPrefix(proxy, constant.Hysteria2Prefix1)
	proxy = strings.TrimPrefix(proxy, constant.Hysteria2Prefix2)
	urlParts := strings.SplitN(proxy, "@", 2)
	if len(urlParts) != 2 {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidStruct,
			Message: "missing character '@' in url",
			Raw:     proxy,
		}
	}
	password := urlParts[0]

	serverInfo := strings.SplitN(urlParts[1], "/?", 2)
	if len(serverInfo) != 2 {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidStruct,
			Message: "missing params in url",
			Raw:     proxy,
		}
	}
	paramStr := serverInfo[1]

	serverAndPort := strings.SplitN(serverInfo[0], ":", 2)
	var server string
	var portStr string
	if len(serverAndPort) == 1 {
		portStr = "443"
	} else if len(serverAndPort) == 2 {
		server, portStr = serverAndPort[0], serverAndPort[1]
	} else {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidStruct,
			Message: "missing server host or port",
			Raw:     proxy,
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

	params, err := url.ParseQuery(paramStr)
	if err != nil {
		return model.Proxy{}, &ParseError{
			Type:    ErrCannotParseParams,
			Raw:     proxy,
			Message: err.Error(),
		}
	}

	remarks, network, obfs, obfsPassword, pinSHA256, insecure, sni := params.Get("name"), params.Get("network"), params.Get("obfs"), params.Get("obfs-password"), params.Get("pinSHA256"), params.Get("insecure"), params.Get("sni")
	enableTLS := pinSHA256 != ""
	insecureBool := insecure == "1"

	if remarks == "" {
		remarksIndex := strings.LastIndex(proxy, "#")
		if remarksIndex != -1 {
			remarks = proxy[remarksIndex:]
			remarks, _ = url.QueryUnescape(remarks)
		}
	}

	result := model.Proxy{
		Type:           "hysteria2",
		Name:           remarks,
		Server:         server,
		Port:           port,
		Password:       password,
		Obfs:           obfs,
		ObfsParam:      obfsPassword,
		Sni:            sni,
		SkipCertVerify: insecureBool,
		TLS:            enableTLS,
		Network:        network,
	}
	return result, nil
}
