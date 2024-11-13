package parser

import (
	"fmt"
	"github.com/nitezs/sub2clash/constant"
	"github.com/nitezs/sub2clash/model"
	"net/url"
	"strings"
)

func ParseSocks(proxy string) (model.Proxy, error) {
	if !strings.HasPrefix(proxy, constant.SocksPrefix) {
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

	remarks := link.Fragment
	if remarks == "" {
		remarks = fmt.Sprintf("%s:%s", server, portStr)
	}
	remarks = strings.TrimSpace(remarks)

	encodeStr := link.User.Username()
	var username, password string
	if encodeStr != "" {
		decodeStr, err := DecodeBase64(encodeStr)
		splitStr := strings.Split(decodeStr, ":")
		if err != nil {
			return model.Proxy{}, &ParseError{
				Type:    ErrInvalidStruct,
				Message: "url parse error",
				Raw:     proxy,
			}
		}
		username = splitStr[0]
		if len(splitStr) == 2 {
			password = splitStr[1]
		}
	}
	return model.Proxy{
		Type:     "socks5",
		Name:     remarks,
		Server:   server,
		Port:     port,
		Username: username,
		Password: password,
	}, nil

}
