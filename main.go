package main

import (
	_ "embed"
	"io"
	"strconv"
	"sub2clash/api"
	"sub2clash/common"
	"sub2clash/common/database"
	"sub2clash/config"
	"sub2clash/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//go:embed templates/template_meta.yaml
var templateMeta string

//go:embed templates/template_clash.yaml
var templateClash string

func init() {
	var err error
	// 创建文件夹
	err = common.MkEssentialDir()
	if err != nil {
		logger.Logger.Panic("create essential dir failed", zap.Error(err))
	}
	// 加载配置
	err = config.LoadConfig()
	// 初始化日志
	logger.InitLogger(config.Default.LogLevel)
	if err != nil {
		logger.Logger.Panic("load config failed", zap.Error(err))
	}
	// 写入默认模板
	err = common.WriteDefalutTemplate(templateMeta, templateClash)
	if err != nil {
		logger.Logger.Panic("write default template failed", zap.Error(err))
	}
	// 连接数据库
	err = database.ConnectDB()
	if err != nil {
		logger.Logger.Panic("database connect failed", zap.Error(err))
	}
	logger.Logger.Info("database connect success")
}

func main() {
	// 设置运行模式
	gin.SetMode(gin.ReleaseMode)
	// 关闭 Gin 的日志输出
	gin.DefaultWriter = io.Discard
	// 创建路由
	r := gin.Default()
	// 设置路由
	api.SetRoute(r)
	logger.Logger.Info("server is running at http://localhost:" + strconv.Itoa(config.Default.Port))
	err := r.Run(":" + strconv.Itoa(config.Default.Port))
	if err != nil {
		logger.Logger.Error("server running failed", zap.Error(err))
		return
	}
}
