package parser

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/nitezs/sub2clash/constant"
	"github.com/nitezs/sub2clash/model"
)

func ParseShadowsocksR(proxy string) (model.Proxy, error) {
	if !strings.HasPrefix(proxy, constant.ShadowsocksRPrefix) {
		return model.Proxy{}, &ParseError{Type: ErrInvalidPrefix, Raw: proxy}
	}

	proxy = strings.TrimPrefix(proxy, constant.ShadowsocksRPrefix)
	proxy, err := DecodeBase64(proxy)
	if err != nil {
		return model.Proxy{}, &ParseError{
			Type: ErrInvalidBase64,
			Raw:  proxy,
		}
	}
	serverInfoAndParams := strings.SplitN(proxy, "/?", 2)
	parts := strings.Split(serverInfoAndParams[0], ":")
	server := parts[0]
	protocol := parts[2]
	method := parts[3]
	obfs := parts[4]
	password, err := DecodeBase64(parts[5])
	if err != nil {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidStruct,
			Raw:     proxy,
			Message: err.Error(),
		}
	}
	port, err := ParsePort(parts[1])
	if err != nil {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidPort,
			Message: err.Error(),
			Raw:     proxy,
		}
	}

	var obfsParam string
	var protoParam string
	var remarks string
	if len(serverInfoAndParams) == 2 {
		params, err := url.ParseQuery(serverInfoAndParams[1])
		if err != nil {
			return model.Proxy{}, &ParseError{
				Type:    ErrCannotParseParams,
				Raw:     proxy,
				Message: err.Error(),
			}
		}
		if params.Get("obfsparam") != "" {
			obfsParam, err = DecodeBase64(params.Get("obfsparam"))
		}
		if params.Get("protoparam") != "" {
			protoParam, err = DecodeBase64(params.Get("protoparam"))
		}
		if params.Get("remarks") != "" {
			remarks, err = DecodeBase64(params.Get("remarks"))
		} else {
			remarks = server + ":" + strconv.Itoa(port)
		}
		if err != nil {
			return model.Proxy{}, &ParseError{
				Type:    ErrInvalidStruct,
				Raw:     proxy,
				Message: err.Error(),
			}
		}
	}

	result := model.Proxy{
		Name:          remarks,
		Type:          "ssr",
		Server:        server,
		Port:          port,
		Protocol:      protocol,
		Cipher:        method,
		Obfs:          obfs,
		Password:      password,
		ObfsParam:     obfsParam,
		ProtocolParam: protoParam,
	}

	return result, nil
}
