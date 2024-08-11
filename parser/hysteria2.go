package parser

import (
	"fmt"
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

	link, err := url.Parse(proxy)
	if err != nil {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidStruct,
			Message: "url parse error",
			Raw:     proxy,
		}
	}

	username := link.User.Username()
	password, exist := link.User.Password()
	if !exist {
		password = username
	}

	query := link.Query()
	server := link.Hostname()
	if server == "" {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidStruct,
			Message: "missing server host",
			Raw:     proxy,
		}
	}
	portStr := link.Port()
	if portStr == "" {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidStruct,
			Message: "missing server port",
			Raw:     proxy,
		}
	}
	port, err := ParsePort(portStr)
	if err != nil {
		return model.Proxy{}, &ParseError{
			Type: ErrInvalidPort,
			Raw:  portStr,
		}
	}
	network, obfs, obfsPassword, pinSHA256, insecure, sni := query.Get("network"), query.Get("obfs"), query.Get("obfs-password"), query.Get("pinSHA256"), query.Get("insecure"), query.Get("sni")
	enableTLS := pinSHA256 != "" || sni != ""
	insecureBool := insecure == "1"
	remarks := link.Fragment
	if remarks == "" {
		remarks = fmt.Sprintf("%s:%s", server, portStr)
	}
	remarks = strings.TrimSpace(remarks)

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
