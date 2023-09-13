package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Port               int
	MetaTemplate       string
	ClashTemplate      string
	RequestRetryTimes  int
	RequestMaxFileSize int64
	CacheExpire        int64
	LogLevel           string
}

var Default *Config

func init() {
	Default = &Config{
		MetaTemplate:       "template_meta.yaml",
		ClashTemplate:      "template_clash.yaml",
		RequestRetryTimes:  3,
		RequestMaxFileSize: 1024 * 1024 * 1,
		Port:               8011,
		CacheExpire:        60 * 5,
		LogLevel:           "info",
	}
	err := godotenv.Load()
	if err != nil {
		return
	}
	if os.Getenv("PORT") != "" {
		atoi, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			panic("PORT invalid")
		}
		Default.Port = atoi
	}
	if os.Getenv("META_TEMPLATE") != "" {
		Default.MetaTemplate = os.Getenv("META_TEMPLATE")
	}
	if os.Getenv("CLASH_TEMPLATE") != "" {
		Default.ClashTemplate = os.Getenv("CLASH_TEMPLATE")
	}
	if os.Getenv("REQUEST_RETRY_TIMES") != "" {
		atoi, err := strconv.Atoi(os.Getenv("REQUEST_RETRY_TIMES"))
		if err != nil {
			panic("REQUEST_RETRY_TIMES invalid")
		}
		Default.RequestRetryTimes = atoi
	}
	if os.Getenv("REQUEST_MAX_FILE_SIZE") != "" {
		atoi, err := strconv.Atoi(os.Getenv("REQUEST_MAX_FILE_SIZE"))
		if err != nil {
			panic("REQUEST_MAX_FILE_SIZE invalid")
		}
		Default.RequestMaxFileSize = int64(atoi)
	}
	if os.Getenv("CACHE_EXPIRE") != "" {
		atoi, err := strconv.Atoi(os.Getenv("CACHE_EXPIRE"))
		if err != nil {
			panic("CACHE_EXPIRE invalid")
		}
		Default.CacheExpire = int64(atoi)
	}
	if os.Getenv("LOG_LEVEL") != "" {
		Default.LogLevel = os.Getenv("LOG_LEVEL")
	}
}
