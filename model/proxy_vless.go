package model

type Vless struct {
	Type              string            `yaml:"type"`
	Name              string            `yaml:"name"`
	Server            string            `yaml:"server"`
	Port              int               `yaml:"port"`
	UUID              string            `yaml:"uuid"`
	Flow              string            `yaml:"flow,omitempty"`
	TLS               bool              `yaml:"tls,omitempty"`
	ALPN              []string          `yaml:"alpn,omitempty"`
	UDP               bool              `yaml:"udp,omitempty"`
	PacketAddr        bool              `yaml:"packet-addr,omitempty"`
	XUDP              bool              `yaml:"xudp,omitempty"`
	PacketEncoding    string            `yaml:"packet-encoding,omitempty"`
	Network           string            `yaml:"network,omitempty"`
	RealityOpts       RealityOptions    `yaml:"reality-opts,omitempty"`
	HTTPOpts          HTTPOptions       `yaml:"http-opts,omitempty"`
	HTTP2Opts         HTTP2Options      `yaml:"h2-opts,omitempty"`
	GrpcOpts          GrpcOptions       `yaml:"grpc-opts,omitempty"`
	WSOpts            WSOptions         `yaml:"ws-opts,omitempty"`
	WSPath            string            `yaml:"ws-path,omitempty"`
	WSHeaders         map[string]string `yaml:"ws-headers,omitempty"`
	SkipCertVerify    bool              `yaml:"skip-cert-verify,omitempty"`
	Fingerprint       string            `yaml:"fingerprint,omitempty"`
	ServerName        string            `yaml:"servername,omitempty"`
	ClientFingerprint string            `yaml:"client-fingerprint,omitempty"`
}

func ProxyToVless(p Proxy) Vless {
	return Vless{
		Type:              "vless",
		Name:              p.Name,
		Server:            p.Server,
		Port:              p.Port,
		UUID:              p.UUID,
		Flow:              p.Flow,
		TLS:               p.TLS,
		ALPN:              p.Alpn,
		UDP:               p.UDP,
		PacketAddr:        p.PacketAddr,
		XUDP:              p.XUDP,
		PacketEncoding:    p.PacketEncoding,
		Network:           p.Network,
		RealityOpts:       p.RealityOpts,
		HTTPOpts:          p.HTTPOpts,
		HTTP2Opts:         p.HTTP2Opts,
		GrpcOpts:          p.GrpcOpts,
		WSOpts:            p.WSOpts,
		WSPath:            p.WSOpts.Path,
		WSHeaders:         p.WSOpts.Headers,
		SkipCertVerify:    p.SkipCertVerify,
		Fingerprint:       p.Fingerprint,
		ServerName:        p.Servername,
		ClientFingerprint: p.ClientFingerprint,
	}
}
