package config

import (
	"fmt"
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
}

var Default *Config

func init() {
	Default = &Config{
		MetaTemplate:       "template_meta.yaml",
		ClashTemplate:      "template_clash.yaml",
		RequestRetryTimes:  3,
		RequestMaxFileSize: 1024 * 1024 * 1,
		Port:               8011,
	}
	err := godotenv.Load()
	if err != nil {
		return
	}
	if os.Getenv("PORT") != "" {
		atoi, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			fmt.Println("PORT 不合法")
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
			fmt.Println("REQUEST_RETRY_TIMES 不合法")
		}
		Default.RequestRetryTimes = atoi
	}
	if os.Getenv("REQUEST_MAX_FILE_SIZE") != "" {
		atoi, err := strconv.Atoi(os.Getenv("REQUEST_MAX_FILE_SIZE"))
		if err != nil {
			fmt.Println("REQUEST_MAX_FILE_SIZE 不合法")
		}
		Default.RequestMaxFileSize = int64(atoi)
	}
}
