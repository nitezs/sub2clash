package model

type ShadowSocksR struct {
	Type          string `yaml:"type"`
	Name          string `yaml:"name"`
	Server        string `yaml:"server"`
	Port          int    `yaml:"port"`
	Password      string `yaml:"password"`
	Cipher        string `yaml:"cipher"`
	Obfs          string `yaml:"obfs"`
	ObfsParam     string `yaml:"obfs-param,omitempty"`
	Protocol      string `yaml:"protocol"`
	ProtocolParam string `yaml:"protocol-param,omitempty"`
	UDP           bool   `yaml:"udp,omitempty"`
}

func ProxyToShadowSocksR(p Proxy) ShadowSocksR {
	return ShadowSocksR{
		Type:          "ssr",
		Name:          p.Name,
		Server:        p.Server,
		Port:          p.Port,
		Password:      p.Password,
		Cipher:        p.Cipher,
		Obfs:          p.Obfs,
		ObfsParam:     p.ObfsParam,
		Protocol:      p.Protocol,
		ProtocolParam: p.ProtocolParam,
		UDP:           p.UDP,
	}
}
