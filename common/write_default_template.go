package common

import (
	"os"
	"path/filepath"
	"sub2clash/config"
)

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
			if file != nil {
				_ = file.Close()
			}
		}(file)
		_, err = file.WriteString(template)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteDefalutTemplate(templateMeta string, templateClash string) error {
	if err := writeTemplate(config.Default.MetaTemplate, templateMeta); err != nil {
		return err
	}
	if err := writeTemplate(config.Default.ClashTemplate, templateClash); err != nil {
		return err
	}
	return nil
}
