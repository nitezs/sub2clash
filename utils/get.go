package utils

import (
	"errors"
	"net/http"
	"sub2clash/config"
	"time"
)

func Get(url string) (resp *http.Response, err error) {
	retryTimes := config.Default.RequestRetryTimes
	haveTried := 0
	retryDelay := time.Second // 延迟1秒再重试
	for haveTried < retryTimes {
		client := &http.Client{}
		client.Timeout = time.Second * 10
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			haveTried++
			time.Sleep(retryDelay)
			continue
		}
		req.Header.Set(
			"User-Agent",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		)
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
