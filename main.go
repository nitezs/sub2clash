package main

import (
	_ "embed"
	"io"
	"strconv"

	"github.com/nitezs/sub2clash/api"
	"github.com/nitezs/sub2clash/common"
	"github.com/nitezs/sub2clash/common/database"
	"github.com/nitezs/sub2clash/config"
	"github.com/nitezs/sub2clash/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//go:embed templates/template_meta.yaml
var templateMeta string

//go:embed templates/template_clash.yaml
var templateClash string

func init() {
	var err error

	err = common.MkEssentialDir()
	if err != nil {
		logger.Logger.Panic("create essential dir failed", zap.Error(err))
	}

	err = config.LoadConfig()

	logger.InitLogger(config.Default.LogLevel)
	if err != nil {
		logger.Logger.Panic("load config failed", zap.Error(err))
	}

	err = common.WriteDefalutTemplate(templateMeta, templateClash)
	if err != nil {
		logger.Logger.Panic("write default template failed", zap.Error(err))
	}

	err = database.ConnectDB()
	if err != nil {
		logger.Logger.Panic("database connect failed", zap.Error(err))
	}
	logger.Logger.Info("database connect success")
}

func main() {

	gin.SetMode(gin.ReleaseMode)

	gin.DefaultWriter = io.Discard

	r := gin.Default()

	api.SetRoute(r)
	logger.Logger.Info("server is running at http://localhost:" + strconv.Itoa(config.Default.Port))
	err := r.Run(":" + strconv.Itoa(config.Default.Port))
	if err != nil {
		logger.Logger.Error("server running failed", zap.Error(err))
		return
	}
}
