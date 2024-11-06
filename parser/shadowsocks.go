package parser

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/nitezs/sub2clash/constant"
	"github.com/nitezs/sub2clash/model"
)

func ParseShadowsocks(proxy string) (model.Proxy, error) {
	if !strings.HasPrefix(proxy, constant.ShadowsocksPrefix) {
		return model.Proxy{}, &ParseError{Type: ErrInvalidPrefix, Raw: proxy}
	}
	if !strings.Contains(proxy, "@") {
		s := strings.SplitN(proxy, "#", 2)
		d, err := DecodeBase64(strings.TrimPrefix(s[0], "ss://"))
		if err != nil {
			return model.Proxy{}, &ParseError{
				Type:    ErrInvalidStruct,
				Message: "url parse error",
				Raw:     proxy,
			}
		}
		if len(s) == 2 {
			proxy = "ss://" + d + "#" + s[1]
		} else {
			proxy = "ss://" + d
		}
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
			Type: ErrInvalidStruct,
			Raw:  proxy,
		}
	}

	method := link.User.Username()
	password, _ := link.User.Password()

	if password == "" {
		user, err := DecodeBase64(method)
		if err == nil {
			methodAndPass := strings.SplitN(user, ":", 2)
			if len(methodAndPass) == 2 {
				method = methodAndPass[0]
				password = methodAndPass[1]
			}
		}
	}
	if isLikelyBase64(password) {
		password, err = DecodeBase64(password)
		if err != nil {
			return model.Proxy{}, &ParseError{
				Type:    ErrInvalidStruct,
				Message: "password decode error",
				Raw:     proxy,
			}
		}
	}

	remarks := link.Fragment
	if remarks == "" {
		remarks = fmt.Sprintf("%s:%s", server, portStr)
	}
	remarks = strings.TrimSpace(remarks)

	result := model.Proxy{
		Type:     "ss",
		Cipher:   method,
		Password: password,
		Server:   server,
		Port:     port,
		Name:     remarks,
	}

	return result, nil
}

func isLikelyBase64(s string) bool {
	if len(s)%4 == 0 && strings.HasSuffix(s, "=") && !strings.Contains(strings.TrimSuffix(s, "="), "=") {
		s = strings.TrimSuffix(s, "=")
		chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
		for _, c := range s {
			if !strings.ContainsRune(chars, c) {
				return false
			}
		}
		return true
	}
	return false
}
