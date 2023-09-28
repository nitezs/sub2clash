package model

type ShadowSocks struct {
	Type              string         `yaml:"type"`
	Name              string         `yaml:"name"`
	Server            string         `yaml:"server"`
	Port              int            `yaml:"port"`
	Password          string         `yaml:"password"`
	Cipher            string         `yaml:"cipher"`
	UDP               bool           `yaml:"udp,omitempty"`
	Plugin            string         `yaml:"plugin,omitempty"`
	PluginOpts        map[string]any `yaml:"plugin-opts,omitempty"`
	UDPOverTCP        bool           `yaml:"udp-over-tcp,omitempty"`
	UDPOverTCPVersion int            `yaml:"udp-over-tcp-version,omitempty"`
	ClientFingerprint string         `yaml:"client-fingerprint,omitempty"`
}

func ProxyToShadowSocks(p Proxy) ShadowSocks {
	return ShadowSocks{
		Type:              "ss",
		Name:              p.Name,
		Server:            p.Server,
		Port:              p.Port,
		Password:          p.Password,
		Cipher:            p.Cipher,
		UDP:               p.UDP,
		Plugin:            p.Plugin,
		PluginOpts:        p.PluginOpts,
		UDPOverTCP:        p.UDPOverTCP,
		UDPOverTCPVersion: p.UDPOverTCPVersion,
		ClientFingerprint: p.ClientFingerprint,
	}
}
