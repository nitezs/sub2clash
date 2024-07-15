package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sub2clash/config"
	"sync"
	"time"
)

var subsDir = "subs"
var fileLock sync.RWMutex

func LoadSubscription(url string, refresh bool, userAgent string) ([]byte, error) {
	if refresh {
		return FetchSubscriptionFromAPI(url, userAgent)
	}
	hash := sha256.Sum224([]byte(url))
	fileName := filepath.Join(subsDir, hex.EncodeToString(hash[:]))
	stat, err := os.Stat(fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return FetchSubscriptionFromAPI(url, userAgent)
	}
	lastGetTime := stat.ModTime().Unix() // 单位是秒
	if lastGetTime+config.Default.CacheExpire > time.Now().Unix() {
		file, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		defer func(file *os.File) {
			if file != nil {
				_ = file.Close()
			}
		}(file)
		fileLock.RLock()
		defer fileLock.RUnlock()
		subContent, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		return subContent, nil
	}
	return FetchSubscriptionFromAPI(url, userAgent)
}

func FetchSubscriptionFromAPI(url string, userAgent string) ([]byte, error) {
	hash := sha256.Sum224([]byte(url))
	fileName := filepath.Join(subsDir, hex.EncodeToString(hash[:]))
	resp, err := Get(url, WithUserAgent(userAgent))
	if err != nil {
		return nil, err
	}
	defer func(resp *http.Response) {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}(resp)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		if file != nil {
			_ = file.Close()
		}
	}(file)
	fileLock.Lock()
	defer fileLock.Unlock()
	_, err = file.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write to sub.yaml: %w", err)
	}
	return data, nil
}
