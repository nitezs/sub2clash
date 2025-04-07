package common

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nitezs/sub2clash/config"
)

type GetConfig struct {
	userAgent string
}

type GetOption func(*GetConfig)

func WithUserAgent(userAgent string) GetOption {
	return func(config *GetConfig) {
		config.userAgent = userAgent
	}
}

func Get(url string, options ...GetOption) (resp *http.Response, err error) {
	retryTimes := config.Default.RequestRetryTimes
	haveTried := 0
	retryDelay := time.Second
	getConfig := GetConfig{}
	for _, option := range options {
		option(&getConfig)
	}
	var req *http.Request
	var get *http.Response
	for haveTried < retryTimes {
		client := &http.Client{}
		//client.Timeout = time.Second * 10
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			haveTried++
			time.Sleep(retryDelay)
			continue
		}
		if getConfig.userAgent != "" {
			req.Header.Set("User-Agent", getConfig.userAgent)
		}
		get, err = client.Do(req)
		if err != nil {
			haveTried++
			time.Sleep(retryDelay)
			continue
		} else {
			if get != nil && get.ContentLength > config.Default.RequestMaxFileSize {
				return nil, errors.New("文件过大")
			}
			return get, nil
		}

	}
	return nil, fmt.Errorf("请求失败：%v", err)
}
