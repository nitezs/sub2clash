package config

import (
	"errors"
	"os"
	"regexp"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               int
	MetaTemplate       string
	ClashTemplate      string
	RequestRetryTimes  int
	RequestMaxFileSize int64
	CacheExpire        int64
	LogLevel           string
	//BasePath           string
	ShortLinkLength int
}

var Default *Config
var Version string
var Dev string

func init() {
	reg := regexp.MustCompile(`^v\d+\.\d+\.\d+$`)
	if reg.MatchString(Version) {
		Dev = "false"
	} else {
		Dev = "true"
	}
}

func LoadConfig() error {
	Default = &Config{
		MetaTemplate:       "template_meta.yaml",
		ClashTemplate:      "template_clash.yaml",
		RequestRetryTimes:  3,
		RequestMaxFileSize: 1024 * 1024 * 1,
		Port:               8011,
		CacheExpire:        60 * 5,
		LogLevel:           "info",
		//BasePath:           "/",
		ShortLinkLength: 6,
	}
	_ = godotenv.Load()
	if os.Getenv("PORT") != "" {
		atoi, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			return errors.New("PORT invalid")
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
			return errors.New("REQUEST_RETRY_TIMES invalid")
		}
		Default.RequestRetryTimes = atoi
	}
	if os.Getenv("REQUEST_MAX_FILE_SIZE") != "" {
		atoi, err := strconv.Atoi(os.Getenv("REQUEST_MAX_FILE_SIZE"))
		if err != nil {
			return errors.New("REQUEST_MAX_FILE_SIZE invalid")
		}
		Default.RequestMaxFileSize = int64(atoi)
	}
	if os.Getenv("CACHE_EXPIRE") != "" {
		atoi, err := strconv.Atoi(os.Getenv("CACHE_EXPIRE"))
		if err != nil {
			return errors.New("CACHE_EXPIRE invalid")
		}
		Default.CacheExpire = int64(atoi)
	}
	if os.Getenv("LOG_LEVEL") != "" {
		Default.LogLevel = os.Getenv("LOG_LEVEL")
	}
	//if os.Getenv("BASE_PATH") != "" {
	//	Default.BasePath = os.Getenv("BASE_PATH")
	//	if Default.BasePath[len(Default.BasePath)-1] != '/' {
	//		Default.BasePath += "/"
	//	}
	//}
	if os.Getenv("SHORT_LINK_LENGTH") != "" {
		atoi, err := strconv.Atoi(os.Getenv("SHORT_LINK_LENGTH"))
		if err != nil {
			return errors.New("SHORT_LINK_LENGTH invalid")
		}
		Default.ShortLinkLength = atoi
	}
	return nil
}
