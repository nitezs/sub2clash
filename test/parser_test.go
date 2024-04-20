package test

import (
	"sub2clash/parser"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestHy2Parser(t *testing.T) {
	res, err := parser.ParseHysteria2("hysteria2://letmein@example.com/?insecure=1&obfs=salamander&obfs-password=gawrgura&pinSHA256=deadbeef&sni=real.example.com")
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	bytes, err := yaml.Marshal(res)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	t.Log(string(bytes))
}
