package test

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	// res, err := parser.ParseTrojan("trojan://Abse64hhjewrs@test.com:8443?type=ws&path=%2Fx&host=test.com&security=tls&fp=&alpn=http%2F1.1&sni=test.com#test")
	// if err != nil {
	// 	t.Log(err.Error())
	// 	t.Fail()
	// }
	// bytes, err := yaml.Marshal(res)
	// if err != nil {
	// 	t.Log(err.Error())
	// 	t.Fail()
	// }
	// t.Log(string(bytes))
	t.Log(strings.SplitN("123456", "/?", 2))

}
