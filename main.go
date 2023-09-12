package main

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"sub/api"
	"sub/config"
)

//go:embed templates/template-meta.yaml
var templateMeta string

//go:embed templates/template-clash.yaml
var templateClash string

func writeTemplate(path string, template string) error {
	tPath := filepath.Join(
		"templates", path,
	)
	if _, err := os.Stat(tPath); os.IsNotExist(err) {
		file, err := os.Create(tPath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(file)
		_, err = file.WriteString(template)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func mkDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func init() {
	if err := mkDir("subs"); err != nil {
		os.Exit(1)
	}
	if err := mkDir("templates"); err != nil {
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
	// 创建路由
	r := gin.Default()
	// 设置路由
	api.SetRoute(r)
	fmt.Println("Server is running at 8011")
	err := r.Run(":8011")
	if err != nil {
		fmt.Println(err)
		return
	}
}
