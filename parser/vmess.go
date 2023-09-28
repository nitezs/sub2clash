package parser

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"sub2clash/model"
)

func ParseVmess(proxy string) (model.Proxy, error) {
	// 判断是否以 vmess:// 开头
	if !strings.HasPrefix(proxy, "vmess://") {
		return model.Proxy{}, errors.New("invalid vmess url")
	}
	// 解码
	base64, err := DecodeBase64(strings.TrimPrefix(proxy, "vmess://"))
	if err != nil {
		return model.Proxy{}, errors.New("invalid vmess url" + err.Error())
	}
	// 解析
	var vmess model.VmessJson
	err = json.Unmarshal([]byte(base64), &vmess)
	if err != nil {
		return model.Proxy{}, errors.New("invalid vmess url" + err.Error())
	}
	// 解析端口
	port := 0
	switch vmess.Port.(type) {
	case string:
		port, err = strconv.Atoi(vmess.Port.(string))
		if err != nil {
			return model.Proxy{}, errors.New("invalid vmess url" + err.Error())
		}
	case float64:
		port = int(vmess.Port.(float64))
	}
	// 解析Aid
	aid := 0
	switch vmess.Aid.(type) {
	case string:
		aid, err = strconv.Atoi(vmess.Aid.(string))
		if err != nil {
			return model.Proxy{}, errors.New("invalid vmess url" + err.Error())
		}
	case float64:
		aid = int(vmess.Aid.(float64))
	}
	// 设置默认值
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
		AlterID:           aid,
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
		result.WSOpts = model.WSOptions{
			Path: vmess.Path,
			Headers: map[string]string{
				"Host": vmess.Host,
			},
		}
	}
	return result, nil
}
