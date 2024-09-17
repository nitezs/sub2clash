package parser

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/nitezs/sub2clash/constant"
	"github.com/nitezs/sub2clash/model"
)

func ParseVless(proxy string) (model.Proxy, error) {
	if !strings.HasPrefix(proxy, constant.VLESSPrefix) {
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
	port, err := ParsePort(portStr)
	if err != nil {
		return model.Proxy{}, &ParseError{
			Type:    ErrInvalidPort,
			Message: err.Error(),
			Raw:     proxy,
		}
	}

	query := link.Query()
	uuid := link.User.Username()
	flow, security, alpnStr, sni, insecure, fp, pbk, sid, path, host, serviceName, _type := query.Get("flow"), query.Get("security"), query.Get("alpn"), query.Get("sni"), query.Get("allowInsecure"), query.Get("fp"), query.Get("pbk"), query.Get("sid"), query.Get("path"), query.Get("host"), query.Get("serviceName"), query.Get("type")

	insecureBool := insecure == "1"
	var alpn []string
	if strings.Contains(alpnStr, ",") {
		alpn = strings.Split(alpnStr, ",")
	} else {
		alpn = nil
	}
	remarks := link.Fragment
	if remarks == "" {
		remarks = fmt.Sprintf("%s:%s", server, portStr)
	}
	remarks = strings.TrimSpace(remarks)

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
