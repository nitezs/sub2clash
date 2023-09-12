package model

type Subscription struct {
	Port               int                     `yaml:"port,omitempty"`
	SocksPort          int                     `yaml:"socks-port,omitempty"`
	AllowLan           bool                    `yaml:"allow-lan,omitempty"`
	Mode               string                  `yaml:"mode,omitempty"`
	LogLevel           string                  `yaml:"log-level,omitempty"`
	ExternalController string                  `yaml:"external-controller,omitempty"`
	Proxies            []Proxy                 `yaml:"proxies,omitempty"`
	ProxyGroups        []ProxyGroup            `yaml:"proxy-groups,omitempty"`
	Rules              []string                `yaml:"rules,omitempty"`
	RuleProviders      map[string]RuleProvider `yaml:"rule-providers,omitempty,omitempty"`
}

type ProxyGroup struct {
	Name          string   `yaml:"name,omitempty"`
	Type          string   `yaml:"type,omitempty"`
	Proxies       []string `yaml:"proxies,omitempty"`
	IsCountryGrop bool     `yaml:"-"`
}

type RuleProvider struct {
	Type     string `yaml:"type,omitempty"`
	Behavior string `yaml:"behavior,omitempty"`
	URL      string `yaml:"url,omitempty"`
	Path     string `yaml:"path,omitempty"`
	Interval int    `yaml:"interval,omitempty"`
}

type Payload struct {
	Rules []string `yaml:"payload,omitempty"`
}
