package utils

import (
	"net/http"
	"time"
)

func GetWithRetry(url string) (resp *http.Response, err error) {
	retryTimes := 3
	haveTried := 0
	retryDelay := time.Second // 延迟1秒再重试
	for haveTried < retryTimes {
		get, err := http.Get(url)
		if err != nil {
			haveTried++
			time.Sleep(retryDelay)
			continue
		} else {
			return get, nil
		}
	}
	return nil, err
}
