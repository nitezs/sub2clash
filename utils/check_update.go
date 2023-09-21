package utils

import (
	"encoding/json"
	"errors"
	"io"
	"sub2clash/config"
	"sub2clash/model"
)

func CheckUpdate() (bool, string, error) {
	get, err := Get("https://api.github.com/repos/nitezs/sub2clash/tags")
	if err != nil {
		return false, "", errors.New("get version info failed" + err.Error())

	}
	var version model.Tags
	all, err := io.ReadAll(get.Body)
	if err != nil {
		return false, "", errors.New("get version info failed" + err.Error())

	}
	err = json.Unmarshal(all, &version)
	if err != nil {
		return false, "", errors.New("get version info failed" + err.Error())

	}
	if version[0].Name == config.Version {
		return false, "", nil
	} else {
		return true, version[0].Name, nil
	}
}
