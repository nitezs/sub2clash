package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sub2clash/model"
)

func ParseVmess(proxy string) (model.Proxy, error) {
	// 判断是否以 vmess:// 开头
	if !strings.HasPrefix(proxy, "vmess://") {
		return model.Proxy{}, fmt.Errorf("无效的 vmess URL")
	}
	// 解码
	base64, err := DecodeBase64(strings.TrimPrefix(proxy, "vmess://"))
	if err != nil {
		return model.Proxy{}, errors.New("无效的 vmess URL")
	}
	// 解析
	var vmess model.Vmess
	err = json.Unmarshal([]byte(base64), &vmess)
	if err != nil {
		return model.Proxy{}, errors.New("无效的 vmess URL")
	}
	// 处理端口
	port, err := strconv.Atoi(strings.TrimSpace(vmess.Port))
	if err != nil {
		return model.Proxy{}, errors.New("无效的 vmess URL")
	}
	if vmess.Scy == "" {
		vmess.Scy = "auto"
	}
	if vmess.Net == "ws" && vmess.Path == "" {
		vmess.Path = "/"
	}
	if vmess.Net == "ws" && vmess.Host == "" {
		vmess.Host = vmess.Add
	}
	// 返回结果
	result := model.Proxy{
		Name:              vmess.Ps,
		Type:              "vmess",
		Server:            vmess.Add,
		Port:              port,
		UUID:              vmess.Id,
		AlterID:           vmess.Aid,
		Cipher:            vmess.Scy,
		UDP:               true,
		TLS:               vmess.Tls == "tls",
		Fingerprint:       vmess.Fp,
		ClientFingerprint: "chrome",
		SkipCertVerify:    true,
		Servername:        vmess.Add,
		Network:           vmess.Net,
	}
	if vmess.Net == "ws" {
		result.WSOpts = model.WSOptsStruct{
			Path: vmess.Path,
			Headers: model.HeaderStruct{
				Host: vmess.Host,
			},
		}
	}
	return result, nil
}
