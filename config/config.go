package config

type Config struct {
	Port          int //TODO: 使用自定义端口
	MetaTemplate  string
	ClashTemplate string
}

var Default *Config

func init() {
	Default = &Config{
		MetaTemplate:  "template-meta.yaml",
		ClashTemplate: "template-clash.yaml",
	}
}
