package model

type Trojan struct {
	Type              string         `yaml:"type"`
	Name              string         `yaml:"name"`
	Server            string         `yaml:"server"`
	Port              int            `yaml:"port"`
	Password          string         `yaml:"password"`
	ALPN              []string       `yaml:"alpn,omitempty"`
	SNI               string         `yaml:"sni,omitempty"`
	SkipCertVerify    bool           `yaml:"skip-cert-verify,omitempty"`
	Fingerprint       string         `yaml:"fingerprint,omitempty"`
	UDP               bool           `yaml:"udp,omitempty"`
	Network           string         `yaml:"network,omitempty"`
	RealityOpts       RealityOptions `yaml:"reality-opts,omitempty"`
	GrpcOpts          GrpcOptions    `yaml:"grpc-opts,omitempty"`
	WSOpts            WSOptions      `yaml:"ws-opts,omitempty"`
	ClientFingerprint string         `yaml:"client-fingerprint,omitempty"`
}

func ProxyToTrojan(p Proxy) Trojan {
	return Trojan{
		Type:              "trojan",
		Name:              p.Name,
		Server:            p.Server,
		Port:              p.Port,
		Password:          p.Password,
		ALPN:              p.Alpn,
		SNI:               p.Sni,
		SkipCertVerify:    p.SkipCertVerify,
		Fingerprint:       p.Fingerprint,
		UDP:               p.UDP,
		Network:           p.Network,
		RealityOpts:       p.RealityOpts,
		GrpcOpts:          p.GrpcOpts,
		WSOpts:            p.WSOpts,
		ClientFingerprint: p.ClientFingerprint,
	}
}
