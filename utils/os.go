package utils

import (
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
