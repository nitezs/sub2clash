package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"net/http"
	neturl "net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"encoding/json"

	"github.com/nitezs/sub2clash/config"
)

var subsDir = "subs"
var fileLock sync.RWMutex

func LoadSubscription(url *string, refresh bool, userAgent string) ([]byte, error) {
	if refresh {
		return FetchSubscriptionFromAPI(url, userAgent)
	}
	hash := sha256.Sum224([]byte(*url))
	fileName := filepath.Join(subsDir, hex.EncodeToString(hash[:]))
	stat, err := os.Stat(fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return FetchSubscriptionFromAPI(url, userAgent)
	}
	lastGetTime := stat.ModTime().Unix()
	if lastGetTime+config.Default.CacheExpire > time.Now().Unix() {
		// 读取缓存的订阅名
		if !strings.Contains(*url, "#") {
			nameFile := filepath.Join(subsDir, "name.json")
			if nameData, err := os.ReadFile(nameFile); err == nil {
				var names map[string]string
				if err := json.Unmarshal(nameData, &names); err == nil {
					if name, ok := names[hex.EncodeToString(hash[:])]; ok {
						*url = *url + "#" + name
					}
				}
			}
		}
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

func FetchSubscriptionFromAPI(url *string, userAgent string) ([]byte, error) {
	hash := sha256.Sum224([]byte(*url))
	fileName := filepath.Join(subsDir, hex.EncodeToString(hash[:]))
	resp, err := Get(*url, WithUserAgent(userAgent))
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

	// 从响应头中获取订阅名称
	if !strings.Contains(*url, "#") {
		subName := ""
		if _, params, _ := mime.ParseMediaType(resp.Header.Get("Content-Disposition")); params != nil {
			if fn, ok := params["filename"]; ok {
				subName = fn
			} else if fn, ok := params["filename*"]; ok && strings.HasPrefix(fn, "UTF-8''") {
				if decodedName, err := neturl.QueryUnescape(strings.TrimPrefix(fn, "UTF-8''")); err == nil {
					subName = decodedName
				}
			}
		}

		// 保存订阅名称到 JSON 文件
		if subName != "" {
			nameFile := filepath.Join(subsDir, "name.json")
			var names map[string]string
			if nameData, err := os.ReadFile(nameFile); err == nil {
				_ = json.Unmarshal(nameData, &names)
			}
			if names == nil {
				names = make(map[string]string)
			}
			hash := sha256.Sum224([]byte(*url))
			names[hex.EncodeToString(hash[:])] = subName
			if nameData, err := json.Marshal(names); err == nil {
				_ = os.WriteFile(nameFile, nameData, 0644)
			}
		}

		*url = *url + "#" + subName
	}
	return data, nil
}
