package parser

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/nitezs/sub2clash/constant"
	"github.com/nitezs/sub2clash/model"
)

func ParseTrojan(proxy string) (model.Proxy, error) {
	if !strings.HasPrefix(proxy, constant.TrojanPrefix) {
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

	password := link.User.Username()
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
			Type:    ErrInvalidPort,
			Message: err.Error(),
			Raw:     proxy,
		}
	}

	remarks := link.Fragment
	if remarks == "" {
		remarks = fmt.Sprintf("%s:%s", server, portStr)
	}
	remarks = strings.TrimSpace(remarks)

	query := link.Query()
	network, security, alpnStr, sni, pbk, sid, fp, path, host, serviceName := query.Get("type"), query.Get("security"), query.Get("alpn"), query.Get("sni"), query.Get("pbk"), query.Get("sid"), query.Get("fp"), query.Get("path"), query.Get("host"), query.Get("serviceName")

	var alpn []string
	if strings.Contains(alpnStr, ",") {
		alpn = strings.Split(alpnStr, ",")
	} else {
		alpn = nil
	}

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

	if network == "grpc" {
		result.GrpcOpts = model.GrpcOptions{
			GrpcServiceName: serviceName,
		}
	}

	return result, nil
}
