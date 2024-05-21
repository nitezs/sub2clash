package model

type SmuxStruct struct {
	Enabled bool `yaml:"enable"`
}

type Proxy struct {
	Name                string                `yaml:"name,omitempty"`
	Server              string                `yaml:"server,omitempty"`
	Port                int                   `yaml:"port,omitempty"`
	Type                string                `yaml:"type,omitempty"`
	Cipher              string                `yaml:"cipher,omitempty"`
	Password            string                `yaml:"password,omitempty"`
	UDP                 bool                  `yaml:"udp,omitempty"`
	UUID                string                `yaml:"uuid,omitempty"`
	Network             string                `yaml:"network,omitempty"`
	Flow                string                `yaml:"flow,omitempty"`
	TLS                 bool                  `yaml:"tls,omitempty"`
	ClientFingerprint   string                `yaml:"client-fingerprint,omitempty"`
	Plugin              string                `yaml:"plugin,omitempty"`
	PluginOpts          map[string]any        `yaml:"plugin-opts,omitempty"`
	Smux                SmuxStruct            `yaml:"smux,omitempty"`
	Sni                 string                `yaml:"sni,omitempty"`
	AllowInsecure       bool                  `yaml:"allow-insecure,omitempty"`
	Fingerprint         string                `yaml:"fingerprint,omitempty"`
	SkipCertVerify      bool                  `yaml:"skip-cert-verify,omitempty"`
	Alpn                []string              `yaml:"alpn,omitempty"`
	XUDP                bool                  `yaml:"xudp,omitempty"`
	Servername          string                `yaml:"servername,omitempty"`
	WSOpts              WSOptions             `yaml:"ws-opts,omitempty"`
	AlterID             int                   `yaml:"alterId,omitempty"`
	GrpcOpts            GrpcOptions           `yaml:"grpc-opts,omitempty"`
	RealityOpts         RealityOptions        `yaml:"reality-opts,omitempty"`
	Protocol            string                `yaml:"protocol,omitempty"`
	Obfs                string                `yaml:"obfs,omitempty"`
	ObfsParam           string                `yaml:"obfs-param,omitempty"`
	ProtocolParam       string                `yaml:"protocol-param,omitempty"`
	Remarks             []string              `yaml:"remarks,omitempty"`
	HTTPOpts            HTTPOptions           `yaml:"http-opts,omitempty"`
	HTTP2Opts           HTTP2Options          `yaml:"h2-opts,omitempty"`
	PacketAddr          bool                  `yaml:"packet-addr,omitempty"`
	PacketEncoding      string                `yaml:"packet-encoding,omitempty"`
	GlobalPadding       bool                  `yaml:"global-padding,omitempty"`
	AuthenticatedLength bool                  `yaml:"authenticated-length,omitempty"`
	UDPOverTCP          bool                  `yaml:"udp-over-tcp,omitempty"`
	UDPOverTCPVersion   int                   `yaml:"udp-over-tcp-version,omitempty"`
	SubName             string                `yaml:"-"`
	Up                  string                `yaml:"up,omitempty"`
	Down                string                `yaml:"down,omitempty"`
	CustomCA            string                `yaml:"ca,omitempty"`
	CustomCAString      string                `yaml:"ca-str,omitempty"`
	CWND                int                   `yaml:"cwnd,omitempty"`
	Auth                string                `yaml:"auth,omitempty"`
	ReceiveWindowConn   int                   `yaml:"recv-window-conn,omitempty"`
	ReceiveWindow       int                   `yaml:"recv-window,omitempty"`
	DisableMTUDiscovery bool                  `yaml:"disable-mtu-discovery,omitempty"`
	FastOpen            bool                  `yaml:"fast-open,omitempty"`
	HopInterval         int                   `yaml:"hop-interval,omitempty"`
	Ports               string                `yaml:"ports,omitempty"`
	AuthStringOLD       string                `yaml:"auth_str,omitempty"`
	AuthString          string                `yaml:"auth-str,omitempty"`
	Ip                  string                `yaml:"ip,omitempty"`
	Ipv6                string                `yaml:"ipv6,omitempty"`
	PrivateKey          string                `yaml:"private-key,omitempty"`
	Workers             int                   `yaml:"workers,omitempty"`
	MTU                 int                   `yaml:"mtu,omitempty"`
	PersistentKeepalive int                   `yaml:"persistent-keepalive,omitempty"`
	Peers               []WireGuardPeerOption `yaml:"peers,omitempty"`
	RemoteDnsResolve    bool                  `yaml:"remote-dns-resolve,omitempty"`
	Dns                 []string              `yaml:"dns,omitempty"`
}

type WireGuardPeerOption struct {
	Server       string   `yaml:"server"`
	Port         int      `yaml:"port"`
	PublicKey    string   `yaml:"public-key,omitempty"`
	PreSharedKey string   `yaml:"pre-shared-key,omitempty"`
	Reserved     []uint8  `yaml:"reserved,omitempty"`
	AllowedIPs   []string `yaml:"allowed-ips,omitempty"`
}

type _Proxy Proxy

func (p Proxy) MarshalYAML() (interface{}, error) {
	switch p.Type {
	case "vmess":
		return ProxyToVmess(p), nil
	case "ss":
		return ProxyToShadowSocks(p), nil
	case "ssr":
		return ProxyToShadowSocksR(p), nil
	case "vless":
		return ProxyToVless(p), nil
	case "trojan":
		return ProxyToTrojan(p), nil
	case "hysteria":
		return ProxyToHysteria(p), nil
	case "hysteria2":
		return ProxyToHysteria2(p), nil
	default:
		return _Proxy(p), nil
	}
}
