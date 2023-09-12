package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LoadTemplate 加载模板
// template 模板文件名
func LoadTemplate(template string) (string, error) {
	tPath := filepath.Join("templates", template)
	if _, err := os.Stat(tPath); err == nil {
		file, err := os.Open(tPath)
		if err != nil {
			return "", err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(file)
		result, err := io.ReadAll(file)
		if err != nil {
			return "", err
		}
		if err != nil {
			return "", err
		}
		return string(result), nil
	}
	return "", errors.New("模板文件不存在")
}
