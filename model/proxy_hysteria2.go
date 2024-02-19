package model

type Hysteria2 struct {
	Type           string   `yaml:"type"`
	Name           string   `yaml:"name"`
	Server         string   `yaml:"server"`
	Port           int      `yaml:"port"`
	Up             string   `yaml:"up,omitempty"`
	Down           string   `yaml:"down,omitempty"`
	Password       string   `yaml:"password,omitempty"`
	Obfs           string   `yaml:"obfs,omitempty"`
	ObfsPassword   string   `yaml:"obfs-password,omitempty"`
	SNI            string   `yaml:"sni,omitempty"`
	SkipCertVerify bool     `yaml:"skip-cert-verify,omitempty"`
	Fingerprint    string   `yaml:"fingerprint,omitempty"`
	ALPN           []string `yaml:"alpn,omitempty"`
	CustomCA       string   `yaml:"ca,omitempty"`
	CustomCAString string   `yaml:"ca-str,omitempty"`
	CWND           int      `yaml:"cwnd,omitempty"`
}

func ProxyToHysteria2(p Proxy) Hysteria2 {
	return Hysteria2{
		Type:           "hysteria2",
		Name:           p.Name,
		Server:         p.Server,
		Port:           p.Port,
		Up:             p.Up,
		Down:           p.Down,
		Password:       p.Password,
		Obfs:           p.Obfs,
		ObfsPassword:   p.ObfsParam,
		SNI:            p.Sni,
		SkipCertVerify: p.SkipCertVerify,
		Fingerprint:    p.Fingerprint,
		ALPN:           p.Alpn,
		CustomCA:       p.CustomCA,
		CustomCAString: p.CustomCAString,
		CWND:           p.CWND,
	}
}
