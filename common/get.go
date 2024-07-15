package common

import (
	"errors"
	"net/http"
	"sub2clash/config"
	"time"
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
	retryDelay := time.Second // 延迟1秒再重试
	getConfig := GetConfig{}
	for _, option := range options {
		option(&getConfig)
	}
	for haveTried < retryTimes {
		client := &http.Client{}
		//client.Timeout = time.Second * 10
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			haveTried++
			time.Sleep(retryDelay)
			continue
		}
		if getConfig.userAgent != "" {
			req.Header.Set("User-Agent", getConfig.userAgent)
		}
		get, err := client.Do(req)
		if err != nil {
			haveTried++
			time.Sleep(retryDelay)
			continue
		} else {
			// 如果文件大小大于设定，直接返回错误
			if get != nil && get.ContentLength > config.Default.RequestMaxFileSize {
				return nil, errors.New("文件过大")
			}
			return get, nil
		}

	}
	return nil, err
}
