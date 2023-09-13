package main

import (
	_ "embed"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sub2clash/api"
	"sub2clash/config"
	"sub2clash/logger"
	"sub2clash/utils"
)

//go:embed templates/template_meta.yaml
var templateMeta string

//go:embed templates/template_clash.yaml
var templateClash string

func writeTemplate(path string, template string) error {
	tPath := filepath.Join(
		"templates", path,
	)
	if _, err := os.Stat(tPath); os.IsNotExist(err) {
		file, err := os.Create(tPath)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		_, err = file.WriteString(template)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	if err := utils.MKDir("subs"); err != nil {
		os.Exit(1)
	}
	if err := utils.MKDir("templates"); err != nil {
		os.Exit(1)
	}
	if err := writeTemplate(config.Default.MetaTemplate, templateMeta); err != nil {
		os.Exit(1)
	}
	if err := writeTemplate(config.Default.ClashTemplate, templateClash); err != nil {
		os.Exit(1)
	}
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
	logger.Logger.Info("Server is running at http://localhost:" + strconv.Itoa(config.Default.Port))
	err := r.Run(":" + strconv.Itoa(config.Default.Port))
	if err != nil {
		logger.Logger.Error("Server run error", zap.Error(err))
		return
	}
}
