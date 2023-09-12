package parser

import (
	"encoding/base64"
)

func DecodeBase64(s string) (string, error) {
	decodeStr, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(decodeStr), nil
}
