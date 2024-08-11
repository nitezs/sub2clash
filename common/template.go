package common

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

func LoadTemplate(template string) ([]byte, error) {
	tPath := filepath.Join("templates", template)
	if _, err := os.Stat(tPath); err == nil {
		file, err := os.Open(tPath)
		if err != nil {
			return nil, err
		}
		defer func(file *os.File) {
			if file != nil {
				_ = file.Close()
			}
		}(file)
		result, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, errors.New("模板文件不存在")
}
