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
	needDecode := true
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
		needDecode = false
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

	method := ""
	password := ""
	if needDecode {
		user, err := DecodeBase64(link.User.Username())
		if err != nil {
			return model.Proxy{}, &ParseError{
				Type:    ErrInvalidStruct,
				Message: "missing method and password",
				Raw:     proxy,
			}
		}
		if user == "" {
			return model.Proxy{}, &ParseError{
				Type:    ErrInvalidStruct,
				Message: "missing method and password",
				Raw:     proxy,
			}
		}
		methodAndPass := strings.SplitN(user, ":", 2)
		if len(methodAndPass) != 2 {
			return model.Proxy{}, &ParseError{
				Type:    ErrInvalidStruct,
				Message: "missing method and password",
				Raw:     proxy,
			}
		}
		method = methodAndPass[0]
		password = methodAndPass[1]
	} else {
		method = link.User.Username()
		password, _ = link.User.Password()
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
