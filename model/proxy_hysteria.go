package model

type Hysteria struct {
	Type                string   `yaml:"type"`
	Name                string   `yaml:"name"`
	Server              string   `yaml:"server"`
	Port                int      `yaml:"port,omitempty"`
	Ports               string   `yaml:"ports,omitempty"`
	Protocol            string   `yaml:"protocol,omitempty"`
	ObfsProtocol        string   `yaml:"obfs-protocol,omitempty"` // compatible with Stash
	Up                  string   `yaml:"up"`
	UpSpeed             int      `yaml:"up-speed,omitempty"` // compatible with Stash
	Down                string   `yaml:"down"`
	DownSpeed           int      `yaml:"down-speed,omitempty"` // compatible with Stash
	Auth                string   `yaml:"auth,omitempty"`
	AuthStringOLD       string   `yaml:"auth_str,omitempty"`
	AuthString          string   `yaml:"auth-str,omitempty"`
	Obfs                string   `yaml:"obfs,omitempty"`
	SNI                 string   `yaml:"sni,omitempty"`
	SkipCertVerify      bool     `yaml:"skip-cert-verify,omitempty"`
	Fingerprint         string   `yaml:"fingerprint,omitempty"`
	ALPN                []string `yaml:"alpn,omitempty"`
	CustomCA            string   `yaml:"ca,omitempty"`
	CustomCAString      string   `yaml:"ca-str,omitempty"`
	ReceiveWindowConn   int      `yaml:"recv-window-conn,omitempty"`
	ReceiveWindow       int      `yaml:"recv-window,omitempty"`
	DisableMTUDiscovery bool     `yaml:"disable-mtu-discovery,omitempty"`
	FastOpen            bool     `yaml:"fast-open,omitempty"`
	HopInterval         int      `yaml:"hop-interval,omitempty"`
}

func ProxyToHysteria(p Proxy) Hysteria {
	return Hysteria{
		Type:                "hysteria",
		Name:                p.Name,
		Server:              p.Server,
		Port:                p.Port,
		Ports:               p.Ports,
		Protocol:            p.Protocol,
		Up:                  p.Up,
		Down:                p.Down,
		Auth:                p.Auth,
		AuthStringOLD:       p.AuthStringOLD,
		AuthString:          p.AuthString,
		Obfs:                p.Obfs,
		SNI:                 p.Sni,
		SkipCertVerify:      p.SkipCertVerify,
		Fingerprint:         p.Fingerprint,
		ALPN:                p.Alpn,
		CustomCA:            p.CustomCA,
		CustomCAString:      p.CustomCAString,
		ReceiveWindowConn:   p.ReceiveWindowConn,
		ReceiveWindow:       p.ReceiveWindow,
		DisableMTUDiscovery: p.DisableMTUDiscovery,
		FastOpen:            p.FastOpen,
		HopInterval:         p.HopInterval,
	}
}
