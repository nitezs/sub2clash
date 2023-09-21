package utils

import (
	"errors"
	"os"
)

func MKDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {

			return err
		}
	}
	return nil
}

func MkEssentialDir() error {
	if err := MKDir("subs"); err != nil {
		return errors.New("create subs dir failed" + err.Error())
	}
	if err := MKDir("templates"); err != nil {
		return errors.New("create templates dir failed" + err.Error())
	}
	if err := MKDir("logs"); err != nil {
		return errors.New("create logs dir failed" + err.Error())
	}
	return nil
}
