package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var subsDir = "subs"
var fileLock sync.RWMutex

func LoadSubscription(url string, refresh bool) ([]byte, error) {
	if refresh {
		return FetchSubscriptionFromAPI(url)
	}
	hash := md5.Sum([]byte(url))
	fileName := filepath.Join(subsDir, hex.EncodeToString(hash[:]))
	const refreshInterval = 500 * 60 // 5分钟
	stat, err := os.Stat(fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return FetchSubscriptionFromAPI(url)
	}
	lastGetTime := stat.ModTime().Unix() // 单位是秒
	if lastGetTime+refreshInterval > time.Now().Unix() {
		file, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Println(err)
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
	return FetchSubscriptionFromAPI(url)
}

func FetchSubscriptionFromAPI(url string) ([]byte, error) {
	hash := md5.Sum([]byte(url))
	fileName := filepath.Join(subsDir, hex.EncodeToString(hash[:]))
	resp, err := Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	fileLock.Lock()
	defer fileLock.Unlock()
	_, err = file.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write to sub.yaml: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}
	return data, nil
}
