package parser

import (
	"encoding/base64"
	"strings"
)

func DecodeBase64(s string) (string, error) {
	s = strings.TrimSpace(s)
	if len(s)%4 != 0 {
		s += strings.Repeat("=", 4-len(s)%4)
	}
	decodeStr, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(decodeStr), nil
}
