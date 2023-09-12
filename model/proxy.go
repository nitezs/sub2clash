package model

type PluginOptsStruct struct {
	Mode string `yaml:"mode"`
}

type SmuxStruct struct {
	Enabled bool `yaml:"enable"`
}

type HeaderStruct struct {
	Host string `yaml:"Host"`
}

type WSOptsStruct struct {
	Path                string       `yaml:"path,omitempty"`
	Headers             HeaderStruct `yaml:"headers,omitempty"`
	MaxEarlyData        int          `yaml:"max-early-data,omitempty"`
	EarlyDataHeaderName string       `yaml:"early-data-header-name,omitempty"`
}

type Vmess struct {
	V    string `json:"v"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port string `json:"port"`
	Id   string `json:"id"`
	Aid  string `json:"aid"`
	Scy  string `json:"scy"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	Tls  string `json:"tls"`
	Sni  string `json:"sni"`
	Alpn string `json:"alpn"`
	Fp   string `json:"fp"`
}

type GRPCOptsStruct struct {
	GRPCServiceName string `yaml:"grpc-service-name,omitempty"`
}

type RealityOptsStruct struct {
	PublicKey string `yaml:"public-key,omitempty"`
	ShortId   string `yaml:"short-id,omitempty"`
}

type Proxy struct {
	Name              string            `yaml:"name,omitempty"`
	Server            string            `yaml:"server,omitempty"`
	Port              int               `yaml:"port,omitempty"`
	Type              string            `yaml:"type,omitempty"`
	Cipher            string            `yaml:"cipher,omitempty"`
	Password          string            `yaml:"password,omitempty"`
	UDP               bool              `yaml:"udp,omitempty"`
	UUID              string            `yaml:"uuid,omitempty"`
	Network           string            `yaml:"network,omitempty"`
	Flow              string            `yaml:"flow,omitempty"`
	TLS               bool              `yaml:"tls,omitempty"`
	ClientFingerprint string            `yaml:"client-fingerprint,omitempty"`
	UdpOverTcp        bool              `yaml:"udp-over-tcp,omitempty"`
	UdpOverTcpVersion string            `yaml:"udp-over-tcp-version,omitempty"`
	Plugin            string            `yaml:"plugin,omitempty"`
	PluginOpts        PluginOptsStruct  `yaml:"plugin-opts,omitempty"`
	Smux              SmuxStruct        `yaml:"smux,omitempty"`
	Sni               string            `yaml:"sni,omitempty"`
	AllowInsecure     bool              `yaml:"allow-insecure,omitempty"`
	Fingerprint       string            `yaml:"fingerprint,omitempty"`
	SkipCertVerify    bool              `yaml:"skip-cert-verify,omitempty"`
	Alpn              []string          `yaml:"alpn,omitempty"`
	XUDP              bool              `yaml:"xudp,omitempty"`
	Servername        string            `yaml:"servername,omitempty"`
	WSOpts            WSOptsStruct      `yaml:"ws-opts,omitempty"`
	AlterID           string            `yaml:"alterId,omitempty"`
	GRPCOpts          GRPCOptsStruct    `yaml:"grpc-opts,omitempty"`
	RealityOpts       RealityOptsStruct `yaml:"reality-opts,omitempty"`
	Protocol          string            `yaml:"protocol,omitempty"`
	Obfs              string            `yaml:"obfs,omitempty"`
	ObfsParam         string            `yaml:"obfs-param,omitempty"`
	ProtocolParam     string            `yaml:"protocol-param,omitempty"`
	Remarks           []string          `yaml:"remarks,omitempty"`
}
